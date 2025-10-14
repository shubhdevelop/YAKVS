package command

import (
	"fmt"
	"strconv"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// ExpireAtCommand handles the EXPIREAT command
type ExpireAtCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewExpireAtCommand creates a new EXPIREAT command instance
func NewExpireAtCommand(cmd *parser.Command, store *store.Store) *ExpireAtCommand {
	return &ExpireAtCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the EXPIREAT command
func (gc *ExpireAtCommand) Execute() string {
	if len(gc.Command.Args) < 1 {
		fmt.Println("Error: EXPIREAT requires 1 argument (key)")
		return "$-1\r"
	}

	key := gc.Command.Args[0]
	ttl, err := strconv.ParseInt(gc.Command.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing TTL:", err)
		return "$-1\r"
	}
	value := gc.Store.SetTTL(key, ttl) 
	
	if value {
		fmt.Println("+OK\r")
	} else {
		fmt.Println(":0\r") // we expect the key to be set successfully
	}
	return "+OK\r\n"
}

// ExpireAtCommandMeta provides metadata for the EXPIREAT command
type ExpireAtCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// ExpireAtMeta returns the command metadata
func ExpireAtMeta() *ExpireAtCommandMeta {
	return &ExpireAtCommandMeta{
		Name:      "EXPIREAT",
		Syntax:    "EXPIREAT key timestamp",
		HelpShort: "EXPIREAT sets the expiration time for the key in args",
		HelpLong: `
EXPIREAT sets the expiration time for the key in args.

The command returns +OK if the expiration time is set, :0 if the key does not exist.
		`,
		Examples: `
>> SET k1 v1
OK
>> EXPIREAT k1 1735689600
$2
v1
>> EXPIREAT k2 1735689600
$-1
		`,
	}
}
