package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// GetCommand handles the GET command
type TtlCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewGetCommand creates a new GET command instance
func NewTtlCommand(cmd *parser.Command, store *store.Store) *TtlCommand {
	return &TtlCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the GET command
func (gc *TtlCommand) Execute() {
	if len(gc.Command.Args) < 1 {
		fmt.Println("Error: TTL requires 1 argument (key)")
		return
	}

	key := gc.Command.Args[0]
	value := gc.Store.GetTTL(key)
	
	if value == -2 {
		fmt.Println(":-2\r")
	} else {
		fmt.Printf(":%d\r\n", value)
	}
}

// GetCommandMeta provides metadata for the GET command
type TtlCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// GetMeta returns the command metadata
func TtlMeta() *TtlCommandMeta {
	return &TtlCommandMeta{
		Name:      "TTL",
		Syntax:    "TTL key",
		HelpShort: "TTL returns the time-to-live for the key in args",
		HelpLong: `
TTL returns the time-to-live for the key in args.

The command returns the time-to-live for the key in args.
		`,
		Examples: `
>> SET k1 v1
OK
>> TTL k1
:3599
>> TTL k2
:-1
		`,
	}	
}
