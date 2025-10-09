package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// DelCommand handles the DEL command
type DelCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewDelCommand creates a new DEL command instance
func NewDelCommand(cmd *parser.Command, store *store.Store) *DelCommand {
	return &DelCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the DEL command
func (dc *DelCommand) Execute() {
	if len(dc.Command.Args) < 1 {
		fmt.Println("Error: DEL requires 1 argument (key)")
		return
	}

	key := dc.Command.Args[0]
	value := dc.Store.GetValue(key)
	
	if value == nil {
		fmt.Println("$-1\r")
	} else {
		fmt.Println("+OK\r")
	}
}

// DelCommandMeta provides metadata for the DEL command
type DelCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// GetMeta returns the command metadata
func DelMeta() *DelCommandMeta {
	return &DelCommandMeta{
		Name:      "DEL",
		Syntax:    "DEL key",
		HelpShort: "DEL returns the value as a string for the key in args",
		HelpLong: `
DEL deletes the key and its associated value from the store.

The command returns +OK if the key is deleted, otherwise -1.
		`,
		Examples: `
>> SET k1 v1
OK
>> DEL k1
+OK
		`,
	}
}
