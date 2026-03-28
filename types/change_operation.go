package types

import (
	"fmt"
	"strings"
)

type ChangeOperation struct {
	// Id is the operation's OpId. ActorIdx is local to the change's actor table
	// (0 = the change's own actor) and is resolved to the document's actor table in ToChange().
	Id           OpId
	Object       ObjectId
	Key          Key
	Insert       bool
	Action       Action
	Predecessors []OpId
}

func (o ChangeOperation) String() string {
	var res strings.Builder
	res.WriteString("Operation {\n")
	res.WriteString(fmt.Sprintf("  \tId: %v\n", o.Id))
	res.WriteString(fmt.Sprintf("  \tObject: %v\n", o.Object))
	res.WriteString(fmt.Sprintf("  \tKey: %v\n", o.Key))
	res.WriteString(fmt.Sprintf("  \tInsert: %v\n", o.Insert))
	res.WriteString(fmt.Sprintf("  \tAction: %v\n", o.Action))
	res.WriteString(fmt.Sprintf("  \tPredecessors: %v\n", o.Predecessors))
	for i, pred := range o.Predecessors {
		res.WriteString(fmt.Sprintf("  \tPredecessors[%d]: %v\n", i, pred))
	}
	res.WriteString("  }")
	return res.String()
}
