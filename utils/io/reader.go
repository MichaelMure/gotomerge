package ioutil

import (
	"fmt"
	"io"
)

type SubReader interface {
	io.Reader

	// SubReader returns an io.Reader for a section of the given size, placed offset bytes after
	// the current position. This does NOT change the position. After being done with one or more
	// sub-readers, the caller should call Skip() to move the position past the now irrelevant bytes.
	SubReader(offset uint64, size uint64) (SubReader, error)

	// SubReaderOffset is the same as SubReader but with only an offset. The end of the data remain the same.
	SubReaderOffset(offset uint64) (SubReader, error)

	// Skip increments the reader position by N bytes and release consumed pages.
	Skip(n int) error

	// Empty tells if the underlying reader has been entirely consumed.
	Empty() bool

	// Consumed returns the number of bytes read, with either Read() or Skip()
	Consumed() int

	// Peek writes into w the n next bytes without moving the internal position.
	Peek(w io.Writer, n int) error
}

const pageSize = 16 * 1024 // 16KB pages

var _ SubReader = &pagedReader{}

// pagedReader perform buffered reads on the source io.Reader while also providing sub-readers
// for concurrent reads in the buffered data. New pages are read from the source automatically as needed.
// The internal position of the reader is incremented by either using Read() or Skip(). As this position
// goes beyond old pages, they get released, but the underlying buffer is kept for further reading.
type pagedReader struct {
	source io.Reader

	pages     [][]byte // ring buffer of page data
	head      int      // index of the first page with unconsumed data
	count     int      // number of pages currently in use
	position  int      // position in bytes from the start of the head page
	available int      // total bytes available for reading in the loaded pages
	consumed  int      // total bytes given to a consumer (through Read(), Skip())
}

// NewPagedReader creates a new paged reader with pre-allocated capacity
func NewPagedReader(source io.Reader) SubReader {
	return &pagedReader{
		source: source,
		pages:  make([][]byte, 8),
	}
}

func (r *pagedReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	// If we don't have enough available data, try to grow
	if r.available < len(p) {
		err := r.grow(len(p))
		if err != nil && err != io.EOF {
			return 0, err
		}
		// Continue even if EOF - we might have some data to return
	}

	totalRead := 0

	for totalRead < len(p) && r.available > 0 {
		availableInPage := len(r.pages[r.head]) - r.position
		wanted := len(p) - totalRead

		if wanted < availableInPage {
			copy(p[totalRead:], r.pages[r.head][r.position:r.position+wanted])
			r.position += wanted
			r.available -= wanted
			totalRead += wanted
			break
		} else {
			copy(p[totalRead:], r.pages[r.head][r.position:])
			r.pages[r.head] = r.pages[r.head][:0] // keep the cap, set len to zero
			r.head = (r.head + 1) % len(r.pages)
			r.count--
			r.position = 0
			r.available -= availableInPage
			totalRead += availableInPage
		}
	}

	if totalRead == 0 {
		return 0, io.EOF
	}

	r.consumed += totalRead

	return totalRead, nil
}

// SubReader returns an io.Reader for a section of the given size, placed offset bytes after
// the current position. This does NOT change the position. After being done with one or more
// sub-readers, the caller should call Skip() to move the position past the now irrelevant bytes.
func (r *pagedReader) SubReader(offset uint64, size uint64) (SubReader, error) {
	return newSubReader(r, offset, size)
}

func (r *pagedReader) SubReaderOffset(offset uint64) (SubReader, error) {
	return newSubReaderOffset(r, offset)
}

// Skip increments the reader position by N bytes and release consumed pages.
func (r *pagedReader) Skip(n int) error {
	if n > r.available {
		err := r.grow(n)
		if err != nil {
			return err
		}
	}
	for n > 0 {
		availableInPage := len(r.pages[r.head]) - r.position
		if n < availableInPage {
			r.position += n
			r.available -= n
			r.consumed += n
			return nil
		} else {
			r.pages[r.head] = r.pages[r.head][:0] // keep the cap, set len to zero
			r.head = (r.head + 1) % len(r.pages)
			r.count--
			r.position = 0
			r.available -= availableInPage
			n -= availableInPage
			r.consumed += availableInPage
		}
	}
	return nil
}

// Empty tells if the underlying reader has been entirely consumed.
func (r *pagedReader) Empty() bool {
	if r.available > 0 {
		return false
	}
	err := r.grow(1)
	return err == io.EOF
}

func (r *pagedReader) Consumed() int {
	return r.consumed
}

func (r *pagedReader) Peek(w io.Writer, n int) error {
	if n == 0 {
		return nil
	}

	// If we don't have enough available data, try to grow
	if r.available < n {
		err := r.grow(n)
		if err != nil {
			return err
		}
	}

	totalRead := 0
	head := r.head
	position := r.position
	available := r.available

	for totalRead < n && r.available > 0 {
		availableInPage := len(r.pages[head]) - position
		wanted := n - totalRead

		if wanted < availableInPage {
			_, err := w.Write(r.pages[head][position : position+wanted])
			if err != nil {
				return err
			}
			position += wanted
			available -= wanted
			totalRead += wanted
			break
		} else {
			_, err := w.Write(r.pages[head][position:])
			if err != nil {
				return err
			}
			head = (head + 1) % len(r.pages)
			position = 0
			available -= availableInPage
			totalRead += availableInPage
		}
	}

	return nil
}

// grow ensures we have at least minBytes available for reading
func (r *pagedReader) grow(minBytes int) error {
	// If we already have enough bytes, nothing to do
	if r.available >= minBytes {
		return nil
	}

	needed := minBytes - r.available

	// Read pages until we have enough bytes
	for needed > 0 {
		n, err := r.loadPage()
		if err != nil {
			return err
		}

		needed -= n

		// If we hit EOF, we're done
		if err == io.EOF {
			break
		}
	}

	return nil
}

// loadPage reads the next page of data, for up to pageSize bytes.
func (r *pagedReader) loadPage() (int, error) {
	if r.count >= len(r.pages) {
		r.growPagesRing()
	}
	tail := (r.head + r.count) % len(r.pages)

	// Get or allocate page at tail position
	if r.pages[tail] == nil {
		r.pages[tail] = make([]byte, pageSize)
	} else {
		// Reset slice length to full page capacity for reuse
		r.pages[tail] = r.pages[tail][:pageSize]
	}

	// Read from source into this page
	n, err := r.source.Read(r.pages[tail])

	// Trim page to actual bytes read
	r.pages[tail] = r.pages[tail][:n]

	if err != nil && err != io.EOF {
		return n, err
	}

	if n == 0 {
		// No more data available
		if err == io.EOF {
			return n, io.EOF
		}
		return n, nil
	}

	// Update counters
	r.count++
	r.available += n

	return n, nil
}

// growPagesRing grows the pages slice when the ring buffer is full
func (r *pagedReader) growPagesRing() {
	newLen := len(r.pages)*2 + 1
	newPages := make([][]byte, newLen)
	for i := 0; i < r.count; i++ {
		newPages[i] = r.pages[(r.head+i)%len(r.pages)]
	}

	// Reset head to 0 since we copied in order
	r.head = 0
	r.pages = newPages
}

var _ SubReader = &subReader{}

type subReader struct {
	r        *pagedReader
	head     int // the page at which the SubReader currently reads
	position int // position in bytes from the start of SubReader's head page
	allowed  int // remaining bytes allowed for reading, similar to an io.LimitedReader
	consumed int // total bytes given to a consumer (through Read(), Skip())
}

func newSubReader(r *pagedReader, offset uint64, size uint64) (*subReader, error) {
	if int(offset)+int(size) > r.available {
		err := r.grow(int(offset) + int(size))
		if err != nil {
			return nil, fmt.Errorf("subReader: grow failed: %w", err)
		}
	}

	head := r.head
	position := r.position
	available := r.available
	toSkip := int(offset)

	for toSkip > 0 && available > 0 {
		availableInPage := len(r.pages[head]) - position

		if toSkip < availableInPage {
			position += toSkip
			available -= toSkip
			toSkip = 0
			break
		} else {
			head = (head + 1) % len(r.pages)
			position = 0
			available -= availableInPage
			toSkip -= availableInPage
		}
	}

	return &subReader{
		r:        r,
		head:     head,
		position: position,
		allowed:  int(size),
	}, nil
}

func (s *subReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	totalRead := 0

	for totalRead < len(p) && s.allowed > 0 {
		availableInPage := len(s.r.pages[s.head]) - s.position
		wanted := len(p) - totalRead

		if wanted < availableInPage || s.allowed < availableInPage {
			// we stay within one page
			if wanted <= s.allowed {
				// not reaching the end of available
				copy(p[totalRead:], s.r.pages[s.head][s.position:s.position+wanted])
				totalRead += wanted
				s.position += wanted
				s.allowed -= wanted
				break
			} else {
				// reaching the end of available
				copy(p[totalRead:], s.r.pages[s.head][s.position:s.position+s.allowed])
				totalRead += s.allowed
				s.position += s.allowed
				s.allowed = 0
				return totalRead, io.EOF
			}
		} else {
			// we move to the next page
			copy(p[totalRead:], s.r.pages[s.head][s.position:s.position+availableInPage])
			totalRead += availableInPage
			s.head = (s.head + 1) % len(s.r.pages)
			s.position = 0
			s.allowed -= availableInPage
		}
	}

	if totalRead == 0 {
		return 0, io.EOF
	}

	s.consumed += totalRead

	return totalRead, nil
}

func (s *subReader) SubReader(offset uint64, size uint64) (SubReader, error) {
	if offset+size > uint64(s.allowed) {
		return nil, fmt.Errorf("subReader: offset + size > available")
	}

	head := s.head
	position := s.position
	available := s.allowed
	toSkip := int(offset)

	for toSkip > 0 && available > 0 {
		availableInPage := len(s.r.pages[head]) - position

		if toSkip < availableInPage {
			position += toSkip
			available -= toSkip
			toSkip = 0
			break
		} else {
			head = (head + 1) % len(s.r.pages)
			position = 0
			available -= availableInPage
			toSkip -= availableInPage
		}
	}

	return &subReader{
		r:        s.r,
		head:     head,
		position: position,
		allowed:  int(size),
	}, nil
}

func (s *subReader) SubReaderOffset(offset uint64) (SubReader, error) {
	if offset > uint64(s.allowed) {
		return nil, fmt.Errorf("subReader: offset > available")
	}

	head := s.head
	position := s.position
	available := s.allowed
	toSkip := int(offset)

	for toSkip > 0 && available > 0 {
		availableInPage := len(s.r.pages[head]) - position

		if toSkip < availableInPage {
			position += toSkip
			available -= toSkip
			toSkip = 0
			break
		} else {
			head = (head + 1) % len(s.r.pages)
			position = 0
			available -= availableInPage
			toSkip -= availableInPage
		}
	}

	return &subReader{
		r:        s.r,
		head:     head,
		position: position,
		allowed:  s.allowed - int(offset),
	}, nil
}

func (s *subReader) Skip(n int) error {
	if n > s.allowed {
		return fmt.Errorf("subReader: skip %d > available %d", n, s.allowed)
	}
	for n > 0 {
		availableInPage := len(s.r.pages[s.head]) - s.position
		if n < availableInPage {
			s.position += n
			s.allowed -= n
			s.consumed += n
			return nil
		} else {
			s.r.pages[s.head] = s.r.pages[s.head][:0] // keep the cap, set len to zero
			s.head = (s.head + 1) % len(s.r.pages)
			s.position = 0
			s.allowed -= availableInPage
			n -= availableInPage
			s.consumed += availableInPage
		}
	}
	return nil
}

func (s *subReader) Empty() bool {
	return s.allowed <= 0
}

func (s *subReader) Consumed() int {
	return s.consumed
}

func (s *subReader) Peek(w io.Writer, n int) error {
	panic("not implemented")
}

var _ SubReader = &subReaderOffset{}

type subReaderOffset struct {
	r        *pagedReader
	head     int // the page at which the SubReader currently reads
	position int // position in bytes from the start of SubReader's head page
	consumed int // total bytes given to a consumer (through Read(), Skip())
}

func newSubReaderOffset(r *pagedReader, offset uint64) (*subReaderOffset, error) {
	// Only grow enough to reach the offset, not more
	if int(offset) > r.available {
		err := r.grow(int(offset))
		if err != nil {
			return nil, fmt.Errorf("subReaderOffset: grow failed: %w", err)
		}
	}

	head := r.head
	position := r.position
	available := r.available
	toSkip := int(offset)

	for toSkip > 0 && available > 0 {
		availableInPage := len(r.pages[head]) - position

		if toSkip <= availableInPage {
			position += toSkip
			available -= toSkip
			toSkip = 0
			break
		} else {
			head = (head + 1) % len(r.pages)
			position = 0
			available -= availableInPage
			toSkip -= availableInPage
		}
	}

	return &subReaderOffset{
		r:        r,
		head:     head,
		position: position,
	}, nil
}

func (s *subReaderOffset) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	// if the current page is empty, try to load
	if s.head == (s.r.head+s.r.count)%len(s.r.pages) && len(s.r.pages[s.head]) == 0 {
		_, err := s.r.loadPage()
		if err != nil && err != io.EOF {
			return 0, err
		}
	}

	totalRead := 0

	for totalRead < len(p) {
		availableInPage := len(s.r.pages[s.head]) - s.position
		wanted := len(p) - totalRead

		if wanted <= availableInPage {
			copy(p[totalRead:], s.r.pages[s.head][s.position:s.position+wanted])
			s.position += wanted
			totalRead += wanted
			break
		} else {
			copy(p[totalRead:], s.r.pages[s.head][s.position:])
			totalRead += availableInPage
			if (s.head+1)%len(s.r.pages) == (s.r.head+s.r.count)%len(s.r.pages) {
				loaded, err := s.r.loadPage()
				if loaded == 0 {
					s.consumed += totalRead
					return totalRead, io.EOF
				}
				if err != nil {
					s.consumed += totalRead
					return totalRead, err
				}
			}
			s.head = (s.head + 1) % len(s.r.pages)
			s.position = 0
		}
	}

	if totalRead == 0 {
		return 0, io.EOF
	}

	s.consumed += totalRead

	return totalRead, nil
}

func (s *subReaderOffset) SubReader(offset uint64, size uint64) (SubReader, error) {
	head := s.head
	position := s.position
	toSkip := int(offset)

	for toSkip > 0 {
		availableInPage := len(s.r.pages[head]) - position

		if toSkip < availableInPage {
			position += toSkip
			toSkip = 0
			break
		} else {
			if (head+1)%len(s.r.pages) == (s.r.head+s.r.count)%len(s.r.pages) {
				_, err := s.r.loadPage()
				if err != nil {
					return nil, fmt.Errorf("subReaderOffset: SubReader offset beyond available data: %w", err)
				}
			}
			head = (head + 1) % len(s.r.pages)
			position = 0
			toSkip -= availableInPage
		}
	}

	return &subReader{
		r:        s.r,
		head:     head,
		position: position,
		allowed:  int(size),
	}, nil
}

func (s *subReaderOffset) SubReaderOffset(offset uint64) (SubReader, error) {
	panic("not implemented")
}

func (s *subReaderOffset) Skip(n int) error {
	// if the current page is empty, try to load
	if s.head == (s.r.head+s.r.count)%len(s.r.pages) && len(s.r.pages[s.head]) == 0 {
		_, err := s.r.loadPage()
		if err != nil && err != io.EOF {
			return err
		}
	}

	for n > 0 {
		availableInPage := len(s.r.pages[s.head]) - s.position
		if n < availableInPage {
			s.position += n
			s.consumed += n
			return nil
		} else {
			n -= availableInPage
			s.consumed += availableInPage
			if (s.head+1)%len(s.r.pages) == (s.r.head+s.r.count)%len(s.r.pages) {
				loaded, err := s.r.loadPage()
				if loaded == 0 {
					return io.EOF
				}
				if err != nil {
					return err
				}
			}
			s.head = (s.head + 1) % len(s.r.pages)
			s.position = 0
		}
	}
	return nil
}

func (s *subReaderOffset) Empty() bool {
	panic("not implemented")
}

func (s *subReaderOffset) Consumed() int {
	return s.consumed
}

func (s *subReaderOffset) Peek(w io.Writer, n int) error {
	if n == 0 {
		return nil
	}

	// if the current page is empty, try to load
	if s.head == (s.r.head+s.r.count)%len(s.r.pages) && len(s.r.pages[s.head])-s.position == 0 {
		_, err := s.r.loadPage()
		if err != nil {
			return err
		}
	}

	head := s.head
	position := s.position

	for n > 0 {
		availableInPage := len(s.r.pages[head]) - position

		if n <= availableInPage {
			_, err := w.Write(s.r.pages[head][position : position+n])
			return err
		} else {
			_, err := w.Write(s.r.pages[head][position:])
			if err != nil {
				return err
			}
			n -= availableInPage
			if (head+1)%len(s.r.pages) == (s.r.head+s.r.count)%len(s.r.pages) {
				loaded, err := s.r.loadPage()
				if loaded == 0 {
					return io.EOF
				}
				if err != nil {
					return err
				}
			}
			head = (head + 1) % len(s.r.pages)
			position = 0
		}
	}

	return nil
}

var _ SubReader = &bytesReader{}

type bytesReader struct {
	data     []byte
	position int
	consumed int // total bytes given to a consumer (through Read(), Skip())
}

func NewBytesReader(data []byte) SubReader {
	return &bytesReader{data: data}
}

func (b *bytesReader) Read(p []byte) (n int, err error) {
	var toRead int
	if len(p) <= len(b.data)-b.position {
		toRead = len(p)
	} else {
		toRead = len(b.data) - b.position
		err = io.EOF
	}
	copy(p, b.data[b.position:b.position+toRead])
	b.position += toRead
	b.consumed += toRead
	return toRead, err
}

func (b *bytesReader) SubReader(offset uint64, size uint64) (SubReader, error) {
	if uint64(b.position)+offset+size > uint64(len(b.data)) {
		return nil, fmt.Errorf("bytesReader: offset + size > len(data)")
	}
	return &bytesReader{data: b.data[uint64(b.position)+offset : uint64(b.position)+offset+size], position: 0}, nil
}

func (b *bytesReader) SubReaderOffset(offset uint64) (SubReader, error) {
	if uint64(b.position)+offset > uint64(len(b.data)) {
		return nil, fmt.Errorf("bytesReader: offset > len(data)")
	}
	return &bytesReader{data: b.data[uint64(b.position)+offset:], position: 0}, nil
}

func (b *bytesReader) Skip(n int) error {
	if b.position+n > len(b.data) {
		return fmt.Errorf("bytesReader: skip %d > len(data)", n)
	}
	b.position += n
	b.consumed += n
	return nil
}

func (b *bytesReader) Empty() bool {
	return b.position >= len(b.data)
}

func (b *bytesReader) Consumed() int {
	return b.consumed
}

func (b *bytesReader) Peek(w io.Writer, n int) error {
	if n > len(b.data)-b.position {
		return fmt.Errorf("bytesReader: peek %d > len(data) - position", n)
	}
	_, err := w.Write(b.data[b.position : b.position+n])
	return err
}
