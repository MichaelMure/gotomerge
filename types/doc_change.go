package types

import (
	"fmt"
	"strings"
)

type DocChange struct {
	ActorId ActorId
	SeqNum  uint64
	MaxOp   uint64
	// Deps contains indices into the document's changes array.
	// These are resolved to []ChangeHash at a higher layer (Change.Deps).
	Deps    []uint64
	Time    Timestamp
	Message *string
	ExtraData any
}

func (dc DocChange) String() string {
	var res strings.Builder
	res.WriteString("DocChange {\n")
	res.WriteString(fmt.Sprintf("  ActorId: %v\n", dc.ActorId))
	res.WriteString(fmt.Sprintf("  SeqNum: %v\n", dc.SeqNum))
	res.WriteString(fmt.Sprintf("  MaxOp: %v\n", dc.MaxOp))
	res.WriteString(fmt.Sprintf("  Deps: %v\n", dc.Deps))
	res.WriteString(fmt.Sprintf("  Time: %v\n", dc.Time))
	res.WriteString(fmt.Sprintf("  Message: %v\n", dc.Message))
	res.WriteString(fmt.Sprintf("  ExtraData: %v\n", dc.ExtraData))
	res.WriteString("  }")
	return res.String()
}
