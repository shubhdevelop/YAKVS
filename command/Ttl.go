package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// TtlCommand handles the TTL command
type TtlCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewTtlCommand creates a new TTL command instance
func NewTtlCommand(cmd *parser.Command, store *store.Store) *TtlCommand {
	return &TtlCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the TTL command
func (gc *TtlCommand) Execute() string {
	if len(gc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'TTL' command\r\n"
	}

	key := gc.Command.Args[0]
	value := gc.Store.GetTTL(key)
	
	if value == -2 {
		return ":-2\r\n"
	} 
	valueStr := fmt.Sprintf("%d", value)
	return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
}

// TtlCommandMeta provides metadata for the TTL command
type TtlCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// TtlCommandMeta returns the command metadata
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
