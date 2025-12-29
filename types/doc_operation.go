package types

import (
	"fmt"
	"strings"
)

type DocOperation struct {
	Object     ObjectId
	Key        Key
	Id         OpId
	Insert     bool
	Action     Action
	Successors []OpId
}

func (o DocOperation) String() string {
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
