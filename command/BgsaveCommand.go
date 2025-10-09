package command

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// BgsaveCommand handles the BGSAVE command
type BgSaveCommand struct {
	Command *parser.Command	
	Store   *store.Store
}

// NewBgsaveCommand creates a new BGSAVE command instance
func NewBgSaveCommand(cmd *parser.Command, store *store.Store) *BgSaveCommand {
	return &BgSaveCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the BGSAVE command
func (dc *BgSaveCommand) Execute() {
	if len(dc.Command.Args) < 1 {
		fmt.Println("Error: BGSAVE doesn't require any arguments")	
		return
	}

	// value := dc.Store.Bgsave()
		
	// if value == nil {
	// 	fmt.Println("$-1\r")
	// } else {
	// 	fmt.Println("+OK\r")
	// }
	fmt.Println("+OK\r")
}

// BgSaveCommandMeta provides metadata for the BGSAVE command
type BgSaveCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// BgSaveMeta returns the command metadata
func BgSaveMeta() *BgSaveCommandMeta {
	return &BgSaveCommandMeta{
		Name:      "BGSAVE",
		Syntax:    "BGSAVE",
		HelpShort: "BGSAVE starts a background save of the database",
		HelpLong: `
BGSAVE starts a background save of the database.

	The command returns +OK if the background save is started, otherwise -1.
		`,
		Examples: `
>> BGSAVE
OK
>> BGSAVE
+OK
		`,
	}
}
