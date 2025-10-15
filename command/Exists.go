package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// GetCommand handles the GET command
type ExistsCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewGetCommand creates a new GET command instance
func NewExistsCommand(cmd *parser.Command, store *store.Store) *ExistsCommand {
	return &ExistsCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the GET command
func (gc *ExistsCommand) Execute() string {
	if len(gc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'EXISTS' command\r\n"
	}

	key := gc.Command.Args[0]
	value := gc.Store.Exists(key)
	
	if value {
		valueStr := fmt.Sprintf("%d", 1)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	} else {
		valueStr := fmt.Sprintf("%d", 0)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	}
}

// GetCommandMeta provides metadata for the GET command
type ExistsCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// GetMeta returns the command metadata
func ExistsMeta() *ExistsCommandMeta {
	return &ExistsCommandMeta{
		Name:      "GET",
		Syntax:    "EXISTS key",
		HelpShort: "EXISTS returns 1 if the key exists, 0 if it doesn't",
		HelpLong: `
EXISTS returns 1 if the key exists, 0 if it doesn't.

The command returns 1 if the key exists, 0 if it doesn't.
		`,
			Examples: `
>> SET k1 v1
OK
>> EXISTS k1
$2
v1
>> EXISTS k2
$-1
		`,
	}
}
