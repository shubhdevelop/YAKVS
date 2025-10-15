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
func (sc *IncreByCommand) Execute() string {
	if len(sc.Command.Args) < 2 {
		return "-ERR wrong number of arguments for 'INCRBY' command\r\n"
	}	

	key := sc.Command.Args[0]
	value := sc.Command.Args[1]

	//change the value to string 
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return "-ERR invalid value: " + err.Error() + "\r\n"
	}
	
	newValue, err := sc.Store.IncreBy(key, valueInt)
	if err != nil {
		return "-ERR " + err.Error() + "\r\n"
	}
	valueStr := fmt.Sprintf("%d", newValue)
	return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
}

// IncreByCommandMeta provides metadata for the INCRBY command
type IncreByCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// IncreByMeta returns the command metadata
func IncreByMeta() *IncreByCommandMeta {
	return &IncreByCommandMeta{
		Name:      "INCRBY",
		Syntax:    "INCRBY key value",
		HelpShort: "INCRBY increments the value for the key in args",
		HelpLong: `
INCRBY increments the value for the key in args.

The command returns the new value if the key is incremented.
		`,
		Examples: `
>> INCRBY k1 1
OK
>> INCRBY k2 2
$2
2
		`,
	}
}
