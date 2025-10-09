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
func (gc *GetCommand) Execute() {
	if len(gc.Command.Args) < 1 {
		fmt.Println("Error: GET requires 1 argument (key)")
		return
	}

	key := gc.Command.Args[0]
	value := gc.Store.GetValue(key)
	
	if value == nil {
		fmt.Println("$-1\r")
	} else {
		valueStr := fmt.Sprintf("%v", value)
		fmt.Printf("$%d\r\n%s\r\n", len(valueStr), valueStr)
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

// GetMeta returns the command metadata
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
