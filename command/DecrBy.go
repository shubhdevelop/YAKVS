package command

import (
	"fmt"
	"strconv"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// DecreByCommand handles the DECRBY command
type DecreByCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewDecreByCommand creates a new DECRBY command instance
func NewDecreByCommand(cmd *parser.Command, store *store.Store) *DecreByCommand {
	return &DecreByCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the DECRBY command
func (sc *DecreByCommand) Execute() {
	if len(sc.Command.Args) < 2 {
		fmt.Println("Error: DECRBY requires 2 arguments (key, value)")
		return
	}	

	key := sc.Command.Args[0]
	value := sc.Command.Args[1]

	//change the value to string 
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Error: DECRBY requires a valid integer value")
		return
	}
	
	newValue, err := sc.Store.DecreBy(key, valueInt)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	
	fmt.Printf(":%d\r\n", newValue)
}

// DecreByCommandMeta provides metadata for the DECRBY command
type DecreByCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// DecreByMeta returns the command metadata
func DecreByMeta() *DecreByCommandMeta {
	return &DecreByCommandMeta{
		Name:      "DECRBY",
		Syntax:    "DECRBY key value",
		HelpShort: "DECRBY decrements the value for the key in args",
		HelpLong: `
DECRBY decrements the value for the key in args.

The command returns +OK if the key is decremented.
		`,
		Examples: `
>> DECRBY k1 1
OK
>> DECRBY k2 2
+OK
		`,
	}
}
