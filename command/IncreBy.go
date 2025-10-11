package command

import (
	"fmt"
	"strconv"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// IncreByCommand handles the INCRBY command
type IncreByCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewIncreByCommand creates a new INCRBY command instance
func NewIncreByCommand(cmd *parser.Command, store *store.Store) *IncreByCommand {
	return &IncreByCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the INCRBY command
func (sc *IncreByCommand) Execute() {
	if len(sc.Command.Args) < 2 {
		fmt.Println("Error: INCRBY requires 2 arguments (key, value)")
		return
	}	

	key := sc.Command.Args[0]
	value := sc.Command.Args[1]

	//change the value to string 
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Error: INCRBY requires a valid integer value")
		return
	}
	
	newValue, err := sc.Store.IncreBy(key, valueInt)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf(":%d\r\n", newValue)
}

// IncreByCommandMeta provides metadata for the INCRBY command
type IncreByCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// SetMeta returns the command metadata
func IncreByMeta() *IncreByCommandMeta {
	return &IncreByCommandMeta{
		Name:      "INCRBY",
		Syntax:    "INCRBY key value",
		HelpShort: "INCRBY increments the value for the key in args",
		HelpLong: `
INCRBY increments the value for the key in args.

The command returns +OK if the key is incremented.
		`,
		Examples: `
>> INCRBY k1 1
OK
>> INCRBY k2 2
+OK
		`,
	}
}
