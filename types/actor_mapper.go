package types

// ActorMapper translates global OpSet actor indices to chunk-local indices for
// encoding. The mapping is stored as a dense slice: table[globalIdx] = localIdx.
type ActorMapper struct {
	table []uint32
	next  uint32
}

// NewActorMapper creates an empty ActorMapper with capacity for size global
// indices. Call Add to register actors in local index order (0, 1, 2, …).
func NewActorMapper(size int) ActorMapper {
	return ActorMapper{table: make([]uint32, size)}
}

// IdentityActorMapper returns an ActorMapper that maps every index to itself,
// useful in tests where no translation is needed.
func IdentityActorMapper(n uint32) ActorMapper {
	t := make([]uint32, n)
	for i := range t {
		t[i] = uint32(i)
	}
	return ActorMapper{table: t, next: n}
}

// Add assigns the next available local index to globalIdx.
// The first Add call maps globalIdx → 0, the second → 1, and so on.
func (m *ActorMapper) Add(globalIdx uint32) {
	m.table[globalIdx] = m.next
	m.next++
}

// Map returns the mapped index for the given source index.
func (m *ActorMapper) Map(idx uint32) uint32 { return m.table[idx] }

// MapOpId returns a copy of id with its ActorIdx translated.
func (m *ActorMapper) MapOpId(id OpId) OpId {
	return OpId{ActorIdx: m.table[id.ActorIdx], Counter: id.Counter}
}

// MapObjectId returns a copy of obj with its ActorIdx translated.
// Root objects are returned unchanged.
func (m *ActorMapper) MapObjectId(obj ObjectId) ObjectId {
	if obj.IsRoot() {
		return obj
	}
	return ObjectId(m.MapOpId(OpId(obj)))
}

// MapKey returns a copy of key with its ActorIdx translated.
// String keys and the head-sentinel KeyOpId (Counter==0) are returned unchanged.
func (m *ActorMapper) MapKey(key Key) Key {
	if k, ok := key.(KeyOpId); ok && k.Counter != 0 {
		return KeyOpId{ActorIdx: m.table[k.ActorIdx], Counter: k.Counter}
	}
	return key
}
