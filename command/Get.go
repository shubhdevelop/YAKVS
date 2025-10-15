package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// GetCommand handles the GET command
type GetCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewGetCommand creates a new GET command instance
func NewGetCommand(cmd *parser.Command, store *store.Store) *GetCommand {
	return &GetCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the GET command
func (gc *GetCommand) Execute() string {
	if len(gc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'GET' command\r\n"
	}

	key := gc.Command.Args[0]
	value := gc.Store.GetValue(key)
	
	if value == nil {
		return "$-1\r\n"
	} else {
		// Convert value to string and use bulk string format
		valueStr := fmt.Sprintf("%v", value)
		return fmt.Sprintf("$%d\r\n%s\r\n", len(valueStr), valueStr)
	}
}

// GetCommandMeta provides metadata for the GET command
type GetCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// GetCommandMeta returns the command metadata
func GetMeta() *GetCommandMeta {
	return &GetCommandMeta{
		Name:      "GET",
		Syntax:    "GET key",
		HelpShort: "GET returns the value as a string for the key in args",
		HelpLong: `
GET returns the value as a string for the key in args.

The command returns an empty string if the key does not exist.
		`,
		Examples: `
>> SET k1 v1
OK
>> GET k1
$2
v1
>> GET k2
$-1
		`,
	}
}
