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
func (sc *DecrCommand) Execute() string {
	if len(sc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'DECR' command\r\n"
	}	

	key := sc.Command.Args[0]

	
	newValue, err := sc.Store.DecreBy(key, 1)
	if err != nil {
		return "-ERR " + err.Error() + "\r\n"
	}
	valueStr := fmt.Sprintf("%d", newValue)
	return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
}

// DecrCommandMeta provides metadata for the DECR command
type DecrCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// DecrMeta returns the command metadata
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
$2
0
>> DECR k2
$2
1
		`,
	}
}
