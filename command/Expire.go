package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// GetCommand handles the GET command
type ExpireCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewGetCommand creates a new GET command instance
func NewExpireCommand(cmd *parser.Command, store *store.Store) *ExpireCommand {
	return &ExpireCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the GET command
func (gc *ExpireCommand) Execute() {
	if len(gc.Command.Args) < 1 {
		fmt.Println("Error: EXPIRE requires 2 arguments (key, ttl)")
		return
	}

	key := gc.Command.Args[0]
	ttl, err := strconv.ParseInt(gc.Command.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing TTL:", err)
		return
	}
	ttl = time.Now().Unix() + ttl
	value := gc.Store.SetTTL(key, ttl) 
	
	if value {
		fmt.Println("+OK\r")
	} else {
		fmt.Println(":0\r") // we expect the key to be set successfully
	}
}

// GetCommandMeta provides metadata for the GET command
type ExpireCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// GetMeta returns the command metadata
func ExpireMeta() *ExpireCommandMeta {
	return &ExpireCommandMeta{
		Name:      "EXPIRE",
		Syntax:    "EXPIRE key ttl",
		HelpShort: "EXPIRE sets the expiration time for the key in args",
		HelpLong: `
EXPIRE sets the expiration time for the key in args.

The command returns +OK if the expiration time is set, :0 if the key does not exist.
		`,
		Examples: `
>> SET k1 v1
OK
>> EXPIRE k1 3600
$2
v1
>> EXPIRE k2 3600
$-1
		`,
	}
}
