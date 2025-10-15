package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// IncrCommand handles the INCR command
type IncrCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewIncrCommand creates a new INCR command instance
func NewIncrCommand(cmd *parser.Command, store *store.Store) *IncrCommand {
	return &IncrCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the INCR command
func (sc *IncrCommand) Execute() string {
	if len(sc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'INCR' command\r\n"
	}	

	key := sc.Command.Args[0]

	
	newValue, err := sc.Store.IncreBy(key, 1)
	if err != nil {
		return "-ERR " + err.Error() + "\r\n"
	}
	valueStr := fmt.Sprintf("%d", newValue)
	return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
}

// IncrCommandMeta provides metadata for the INCR command
type IncrCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// IncrMeta returns the command metadata
func IncrMeta() *IncrCommandMeta {
	return &IncrCommandMeta{
		Name:      "INCR",
		Syntax:    "INCR key",
		HelpShort: "INCR increments the value for the key in args",
		HelpLong: `
INCR increments the value for the key in args.

	The command returns the new value if the key is incremented.
		`,
		Examples: `
>> INCR k1
$2
1
>> INCR k2
$2
2
		`,
	}
}
