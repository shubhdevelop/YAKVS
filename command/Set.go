package command

import (
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// SetCommand handles the SET command
type SetCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewSetCommand creates a new SET command instance
func NewSetCommand(cmd *parser.Command, store *store.Store) *SetCommand {
	return &SetCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the SET command
func (sc *SetCommand) Execute() string {
	if len(sc.Command.Args) < 2 {
		return "-ERR wrong number of arguments for 'SET' command\r\n"
	}	

	key := sc.Command.Args[0]
	value := sc.Command.Args[1]
	
	sc.Store.SetValue(key, value)
	return "+OK\r\n" 
}

// SetCommandMeta provides metadata for the SET command
type SetCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// SetCommandMeta returns the command metadata
func SetMeta() *SetCommandMeta {
	return &SetCommandMeta{
		Name:      "SET",
		Syntax:    "SET key value",
		HelpShort: "SET sets the value for the key in args",
		HelpLong: `
SET sets the value for the key in args.

The command returns +OK if the key is set.
		`,
		Examples: `
>> SET k1 v1
OK
>> SET k2 v2
+OK
		`,
	}
}
