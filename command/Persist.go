package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// PersistCommand handles the PERSIST command
type PersistCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewPersistCommand creates a new PERSIST command instance
func NewPersistCommand(cmd *parser.Command, store *store.Store) *PersistCommand {
	return &PersistCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the PERSIST command
func (gc *PersistCommand) Execute() string {
	if len(gc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'PERSIST' command\r\n"
	}

	key := gc.Command.Args[0]
	value := gc.Store.RemoveExpiry(key)
	
	if value {
		valueStr := fmt.Sprintf("%d", 1)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	} else {
		valueStr := fmt.Sprintf("%d", 0)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	}
}

// PersistCommandMeta provides metadata for the PERSIST command
type PersistCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// PersistCommandMeta returns the command metadata
func PersistMeta() *PersistCommandMeta {
	return &PersistCommandMeta{
		Name:      "PERSIST",
		Syntax:    "PERSIST key",
		HelpShort: "PERSIST removes the expiration time for the key in args",
		HelpLong: `
PERSIST removes the expiration time for the key in args.

The command returns +OK if the expiration time is removed, :0 if the key does not exist.
		`,
		Examples: `
>> SET k1 v1
OK
>> PERSIST k1
$2
v1
>> PERSIST k2
$-1
		`,
	}
}
