package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

// ExpireCommand handles the EXPIRE command
type ExpireCommand struct {
	Command *parser.Command
	Store   *store.Store
}

// NewExpireCommand creates a new EXPIRE command instance
func NewExpireCommand(cmd *parser.Command, store *store.Store) *ExpireCommand {
	return &ExpireCommand{
		Command: cmd,
		Store:   store,
	}
}

// Execute executes the EXPIRE command
func (gc *ExpireCommand) Execute() string {
	if len(gc.Command.Args) < 1 {
		return "-ERR wrong number of arguments for 'EXPIRE' command\r\n"
	}

	key := gc.Command.Args[0]
	ttl, err := strconv.ParseInt(gc.Command.Args[1], 10, 64)
	if err != nil {
		return "-ERR invalid TTL: " + err.Error() + "\r\n"
	}
	ttl = time.Now().Unix() + ttl
	value := gc.Store.SetTTL(key, ttl) 
	
	if value {
		valueStr := fmt.Sprintf("%d", 1)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)		
	} else {
		valueStr := fmt.Sprintf("%d", 0)
		return fmt.Sprintf("%d\r\n%s\r\n", len(valueStr), valueStr)
	}
}

// ExpireCommandMeta provides metadata for the EXPIRE command
type ExpireCommandMeta struct {
	Name      string
	Syntax    string
	HelpShort string
	HelpLong  string
	Examples  string
}

// ExpireMeta returns the command metadata
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
