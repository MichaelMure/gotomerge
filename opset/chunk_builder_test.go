package opset

import (
	"github.com/jcalabro/leb128"

	"gotomerge/format"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// rleUint64 encodes a sequence of uint64 values as a single literal RLE run.
func rleUint64(vals ...uint64) []byte {
	var b []byte
	b = append(b, leb128.EncodeS64(int64(-len(vals)))...)
	for _, v := range vals {
		b = append(b, leb128.EncodeU64(v)...)
	}
	return b
}

// rleInt64 encodes a sequence of int64 values as a single literal RLE run.
func rleInt64(vals ...int64) []byte {
	var b []byte
	b = append(b, leb128.EncodeS64(int64(-len(vals)))...)
	for _, v := range vals {
		b = append(b, leb128.EncodeS64(v)...)
	}
	return b
}

// rleString encodes a sequence of strings as a single literal RLE run.
func rleString(vals ...string) []byte {
	var b []byte
	b = append(b, leb128.EncodeS64(int64(-len(vals)))...)
	for _, s := range vals {
		b = append(b, leb128.EncodeU64(uint64(len(s)))...)
		b = append(b, s...)
	}
	return b
}

func sub(data []byte) ioutil.SubReader {
	return ioutil.NewBytesReader(data)
}

// opSetKey encodes a value metadata entry for a string of the given length.
// ValueMetadata = (length << 4) | type; string type = 6.
func valueMetaString(length int) uint64 {
	return uint64(length)<<4 | 6
}

// changeChunkSetKey builds a minimal ChangeChunk with one op: set root[key] = value (string).
func changeChunkSetKey(actor types.ActorId, seqNum, startOp uint64, deps []types.ChangeHash, hash types.ChangeHash, key, value string) *format.ChangeChunk {
	return &format.ChangeChunk{
		Hash:         hash,
		Dependencies: deps,
		Actor:        actor,
		SeqNum:       seqNum,
		StartOp:      startOp,
		OpColumns: format.OperationColumns{
			KeyString:     sub(rleString(key)),
			Action:        sub(rleUint64(uint64(types.ActionSet))),
			ValueMetadata: sub(rleUint64(valueMetaString(len(value)))),
			Value:         sub([]byte(value)),
		},
	}
}

// changeChunkDelete builds a minimal ChangeChunk with one op: delete the op
// identified by predActorIdx/predCounter (in the change's local actor space).
func changeChunkDelete(actor types.ActorId, seqNum, startOp uint64, deps []types.ChangeHash, hash types.ChangeHash, key string, predActorIdx uint32, predCounter int64) *format.ChangeChunk {
	return &format.ChangeChunk{
		Hash:         hash,
		Dependencies: deps,
		Actor:        actor,
		SeqNum:       seqNum,
		StartOp:      startOp,
		OpColumns: format.OperationColumns{
			KeyString:          sub(rleString(key)),
			Action:             sub(rleUint64(uint64(types.ActionDelete))),
			PredecessorGroup:   sub(rleUint64(1)),
			PredecessorActorId: sub(rleUint64(uint64(predActorIdx))),
			PredecessorCounter: sub(rleInt64(predCounter)),
		},
	}
}
