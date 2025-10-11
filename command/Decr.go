package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// DecrCommand handles the DECR command
type DecrCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewDecrCommand creates a new DECR command instance
func NewDecrCommand(cmd *parser.Command, store *store.Store) *DecrCommand {
	return &DecrCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the DECR command
func (sc *DecrCommand) Execute() {
	if len(sc.Command.Args) < 1 {
		fmt.Println("Error: DECR requires 1 argument (key)")
		return
	}	

	key := sc.Command.Args[0]

	
	newValue, err := sc.Store.DecreBy(key, 1)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf(":%d\r\n", newValue)
}

// DecrCommandMeta provides metadata for the DECR command
type DecrCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// SetMeta returns the command metadata
func DecrMeta() *DecrCommandMeta {
	return &DecrCommandMeta{
		Name:      "DECR",
		Syntax:    "DECR key",
		HelpShort: "DECR decrements the value for the key in args",
		HelpLong: `
DECR decrements the value for the key in args.

	The command returns the new value if the key is decremented.
		`,
		Examples: `
>> DECR k1
:0
>> DECR k2
:1
		`,
	}
}
