package gotomerge

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"strings"

	"github.com/jcalabro/leb128"

	"gotomerge/column"
	"gotomerge/column/rle"
	"gotomerge/types"
)

type ChangeChunk struct {
	Dependencies []types.ChangeHash
	Actor        types.ActorId
	SeqNum       uint64
	StartOp      uint64
	Time         types.Timestamp
	Message      string
	OtherActors  []types.ActorId

	OperationMetadata column.Metadata
	Operations        [][]any

	ExtraBytes []byte
}

func (c ChangeChunk) String() string {
	var res strings.Builder
	res.WriteString("ChangeChunk {\n")
	res.WriteString(fmt.Sprintf("  Dependencies: %v\n", c.Dependencies))
	res.WriteString(fmt.Sprintf("  Actor: %v\n", c.Actor))
	res.WriteString(fmt.Sprintf("  SeqNum: %v\n", c.SeqNum))
	res.WriteString(fmt.Sprintf("  StartOp: %v\n", c.StartOp))
	res.WriteString(fmt.Sprintf("  Time: %v\n", c.Time))
	res.WriteString(fmt.Sprintf("  Message: %v\n", c.Message))
	res.WriteString(fmt.Sprintf("  OtherActors: %v\n", c.OtherActors))
	for i, metadatum := range c.OperationMetadata {
		res.WriteString(fmt.Sprintf("  OperationMetadata[%d]: %v\n", i, metadatum))
		res.WriteString(fmt.Sprintf("    Values: %v\n", c.Operations[i]))
	}
	res.WriteString(fmt.Sprintf("  ExtraBytes: %v\n", c.ExtraBytes))
	res.WriteString("}\n")
	return res.String()
}

func readChangeChunk(r io.Reader) (*ChangeChunk, error) {
	var res ChangeChunk
	var err error

	res.Dependencies, err = readChangeHashes(r)
	if err != nil {
		return nil, fmt.Errorf("error reading dependencies: %w", err)
	}

	res.Actor, err = types.ReadLengthEncodedActorId(r)
	if err != nil {
		return nil, fmt.Errorf("error reading actor id: %w", err)
	}

	res.SeqNum, err = leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading seq num: %w", err)
	}

	res.StartOp, err = leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading start op: %w", err)
	}

	time, err := leb128.DecodeS64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading time: %w", err)
	}
	res.Time = types.Timestamp(time)

	msgLen, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %w", err)
	}
	res.Message, err = column.ReadStringWithLimitedPrealloc(r, msgLen)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	res.OtherActors, err = readActorIds(r)
	if err != nil {
		return nil, fmt.Errorf("error reading other actors: %w", err)
	}

	var prevValueMetadata []column.ValueMetadata

	res.Operations = make([][]any, len(res.OperationMetadata))
	for i, metadatum := range res.OperationMetadata {
		rCol := io.LimitReader(r, int64(metadatum.Length))
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rCol)
		}

		switch metadatum.Spec.Type() {
		case column.TypeGroup:
			res.Operations[i] = acc(rle.ReadUint64RLE(rCol))
		case column.TypeActor:
			res.Operations[i] = acc(rle.ReadUint64RLE(rCol))
		case column.TypeULEB128:
			res.Operations[i] = acc(column.ReadUlebColumn(rCol))
		case column.TypeDelta:
			res.Operations[i] = acc(column.ReadDeltaColumn(rCol))
		case column.TypeBool:
			res.Operations[i] = acc(column.ReadBooleanColumn(rCol))
		case column.TypeString:
			res.Operations[i] = acc(column.ReadStringColumn(rCol))
		case column.TypeValueMetadata:
			// TODO: HACK just for early visualisation
			buf := &bytes.Buffer{}
			tee := io.TeeReader(rCol, buf)
			res.Operations[i] = acc(column.ReadValueMetadataColumn(tee))
			prevValueMetadata = nil
			for vm, err := range column.ReadValueMetadataColumn(buf) {
				if err != nil {
					panic(err)
				}
				prevValueMetadata = append(prevValueMetadata, vm)
			}
		case column.TypeValue:
			// skip(rCol, metadatum.Spec, metadatum.Length)

			// TODO: HACK just for early visualisation
			res.Operations[i] = acc(column.ReadValueColumn(r, prevValueMetadata))
		}
	}

	if lr, ok := r.(*io.LimitedReader); ok {
		// Don't try to read the extra bytes if we know there is none.
		// This avoids an allocation in io.ReadAll(), as we know that virtually
		// all the changes we read don't have those extra bytes.
		if lr.N == 0 {
			return &res, nil
		}
	}

	res.ExtraBytes, err = io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading extra bytes: %w", err)
	}

	return &res, nil
}
