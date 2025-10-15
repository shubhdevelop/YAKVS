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
func (dc *DelCommand) Execute() string {
	if len(dc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'DEL' command\r\n"
	}

	key := dc.Command.Args[0]
	
	// Check if key exists before attempting to delete
	if !dc.Store.Exists(key) {
		valueStr := fmt.Sprintf("%d", 0)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	}
	
	// Actually delete the key
	deleted := dc.Store.DeleteValue(key)
	if deleted {
		valueStr := fmt.Sprintf("%d", 1)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	} else {
		valueStr := fmt.Sprintf("%d", 0)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
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

// DelMeta returns the command metadata
func DelMeta() *DelCommandMeta {
	return &DelCommandMeta{
		Name:      "DEL",
		Syntax:    "DEL key",
		HelpShort: "DEL returns the value as a string for the key in args",
		HelpLong: `
DEL deletes the key and its associated value from the store.

The command returns the number of keys deleted, which is 1 if the key is deleted, otherwise 0.
		`,
		Examples: `
>> SET k1 v1
OK
>> DEL k1
$2
1
		`,
	}
}
