package format

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// TestChangeChunkRoundTrip reads each listed testdata file, re-encodes every
// change chunk, and verifies the output is byte-for-byte identical to the
// original. This ensures exact wire-format compatibility for the change chunk
// write path.
func TestChangeChunkRoundTrip(t *testing.T) {
	cases := []string{
		"counter_value_is_ok.automerge",
		// Synthetic: val_meta claims 2 bytes for counter 16, but canonical encoding is 1 byte.
		// We correctly decode and re-encode canonically; exact bytes cannot be preserved.
		// "counter_value_has_incorrect_meta.automerge",
		"two_change_chunks.automerge",
		// Synthetic: change payload < 256 bytes but stored compressed; real encoders never do this.
		// We correctly decompress on read and write uncompressed (DEFLATE_MIN_SIZE = 256).
		// "two_change_chunks_compressed.automerge",
		"two_change_chunks_out_of_order.automerge",
		"text-edits.amrg",
	}

	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			if name == "text-edits.amrg" && testing.Short() {
				t.Skip("skipping large file under -short")
			}

			original, err := os.ReadFile("../testdata/" + name)
			require.NoError(t, err)

			r := ioutil.NewSubReader(original)
			var out bytes.Buffer

			for !r.Empty() {
				chunk, toSkip, err := ReadChunk(r)
				require.NoError(t, err)
				require.NoError(t, r.Skip(toSkip))

				cc, ok := chunk.(*ChangeChunk)
				require.True(t, ok, "expected only ChangeChunks in %s", name)

				ops := NewChangeOpsWriter()
				// Operations already use local actor indices (0 = own actor,
				// 1..N = OtherActors). Identity mapper preserves them as-is.
				n := uint32(1 + len(cc.OtherActors))
				mapper := types.IdentityActorMapper(n)
				for op, err := range cc.Operations() {
					require.NoError(t, err)
					ops.Append(op.Object, op.Key, op.Insert, op.Action, op.Predecessors, mapper)
				}
				cc2 := &ChangeChunk{
					Dependencies: cc.Dependencies,
					Actor:        cc.Actor,
					SeqNum:       cc.SeqNum,
					StartOp:      cc.StartOp,
					Time:         cc.Time,
					Message:      cc.Message,
					OtherActors:  cc.OtherActors,
				}
				require.NoError(t, WriteChange(&out, cc2, ops))
			}

			if !bytes.Equal(original, out.Bytes()) {
				printChunks := func(data []byte) string {
					r := ioutil.NewSubReader(data)
					var res strings.Builder
					for !r.Empty() {
						chunk, toSkip, err := ReadChunk(r)
						if err != nil {
							t.Logf("  parse error: %v", err)
							break
						}
						res.WriteString(fmt.Sprint(chunk))
						_ = r.Skip(toSkip)
					}
					return res.String()
				}
				require.Equal(t, printChunks(original), printChunks(out.Bytes()))
			}
			require.Equal(t, original, out.Bytes())
		})
	}
}

// BenchmarkWriteChange encodes a change chunk with varying numbers of
// set-on-root operations, analogous to Rust's "map save" benchmarks.
func BenchmarkWriteChange(b *testing.B) {
	actor := types.ActorId{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	}
	m := types.IdentityActorMapper(1)

	for _, n := range []int{100, 1_000, 10_000} {
		b.Run(fmt.Sprintf("ops=%d", n), func(b *testing.B) {
			keys := make([]types.Key, n)
			for j := 0; j < n; j++ {
				keys[j] = types.KeyString(fmt.Sprintf("key%d", j))
			}
			var lastLen int
			for i := 0; i < b.N; i++ {
				ops := NewChangeOpsWriter()
				for j := 0; j < n; j++ {
					ops.Append(
						types.RootObjectId(),
						keys[j],
						false,
						types.Action{Kind: types.ActionSet, Value: int64(j)},
						nil,
						m,
					)
				}
				if err := ops.flush(); err != nil {
					b.Fatal(err)
				}
				cc := &ChangeChunk{
					Actor:   actor,
					SeqNum:  1,
					StartOp: 1,
					Time:    types.Timestamp(0),
				}
				var buf bytes.Buffer
				if err := WriteChange(&buf, cc, ops); err != nil {
					b.Fatal(err)
				}
				lastLen = buf.Len()
			}
			b.SetBytes(int64(lastLen))
		})
	}
}
