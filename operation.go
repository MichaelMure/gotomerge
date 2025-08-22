package gotomerge

import (
	"fmt"
	"strings"

	"gotomerge/types"
)

type Operation struct {
	Object     types.ObjectId
	Key        types.Key
	Id         types.ObjectId
	Insert     bool
	Action     types.Action
	Successors []Operation
}

func (o Operation) String() string {
	var res strings.Builder
	res.WriteString("Operation {\n")
	res.WriteString(fmt.Sprintf("  \tObject: %v\n", o.Object))
	res.WriteString(fmt.Sprintf("  \tKey: %v\n", o.Key))
	res.WriteString(fmt.Sprintf("  \tId: %v\n", o.Id))
	res.WriteString(fmt.Sprintf("  \tInsert: %v\n", o.Insert))
	res.WriteString(fmt.Sprintf("  \tAction: %v\n", o.Action))
	res.WriteString(fmt.Sprintf("  \tSuccessors: %v\n", o.Successors))
	for i, succ := range o.Successors {
		res.WriteString(fmt.Sprintf("  \tSuccessors[%d]: %v\n", i, succ))
	}
	res.WriteString("  }")
	return res.String()
}
