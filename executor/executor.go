package executor

import (
	"fmt"
	"strings"

	"github.com/shubhdevelop/YAKVS/command"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

func ExecuteCommand(cmd *parser.Command, store *store.Store) (string, error) {
	fmt.Println("Executing command:", cmd)
	switch strings.ToUpper(cmd.Name) {
	case "BGSAVE":
		bgSaveCmd := command.NewBgSaveCommand(cmd, store)
		return bgSaveCmd.Execute(), nil
	case "SET":
		setCmd := command.NewSetCommand(cmd, store)
		return setCmd.Execute(), nil
	case "GET":
		getCmd := command.NewGetCommand(cmd, store)
		return getCmd.Execute(), nil
	case "DEL":
		delCmd := command.NewDelCommand(cmd, store)
		return delCmd.Execute(), nil
	case "EXISTS":
		existsCmd := command.NewExistsCommand(cmd, store)
		return existsCmd.Execute(), nil
	case "TTL":
		ttlCmd := command.NewTtlCommand(cmd, store)
		return ttlCmd.Execute(), nil
	case "EXPIRE":
		expireCmd := command.NewExpireCommand(cmd, store)
		return expireCmd.Execute(), nil
	case "EXPIREAT":
		expireAtCmd := command.NewExpireAtCommand(cmd, store)
		return expireAtCmd.Execute(), nil
	case "PERSIST":
		persistCmd := command.NewPersistCommand(cmd, store)
		return persistCmd.Execute(), nil
	case "INCRBY":
		incrByCmd := command.NewIncreByCommand(cmd, store)
		return incrByCmd.Execute(), nil
	case "DECRBY":
		decrByCmd := command.NewDecreByCommand(cmd, store)
		return decrByCmd.Execute(), nil	
	case "INCR":
		incrByCmd := command.NewIncrCommand(cmd, store)
		return incrByCmd.Execute(), nil
	case "DECR":
		decrByCmd := command.NewDecrCommand(cmd, store)
		return decrByCmd.Execute(), nil
	default:
		return "", fmt.Errorf("invalid command: %s", cmd.Name)
	}
}




func ExecuteCommandIntegration(cmd *parser.Command, store *store.Store) {
	fmt.Println("Executing command:", cmd)
	switch strings.ToUpper(cmd.Name) {

	case "SET":
		setCmd := command.NewSetCommand(cmd, store)
		setCmd.Execute()
	case "GET":
		getCmd := command.NewGetCommand(cmd, store)
		getCmd.Execute()
	case "DEL":
		delCmd := command.NewDelCommand(cmd, store)
		delCmd.Execute()
	case "EXISTS":
		existsCmd := command.NewExistsCommand(cmd, store)
		existsCmd.Execute()
	case "TTL":
		ttlCmd := command.NewTtlCommand(cmd, store)
		ttlCmd.Execute()
	case "EXPIRE":
		expireCmd := command.NewExpireCommand(cmd, store)
		expireCmd.Execute()
	case "EXPIREAT":
		expireAtCmd := command.NewExpireAtCommand(cmd, store)
		expireAtCmd.Execute()
	case "PERSIST":
		persistCmd := command.NewPersistCommand(cmd, store)
		persistCmd.Execute()
	case "INCRBY":
		incrByCmd := command.NewIncreByCommand(cmd, store)
		incrByCmd.Execute()
	case "DECRBY":
		decrByCmd := command.NewDecreByCommand(cmd, store)
		decrByCmd.Execute()
	case "INCR":
		incrByCmd := command.NewIncrCommand(cmd, store)
		incrByCmd.Execute()
	case "DECR":
		decrByCmd := command.NewDecrCommand(cmd, store)
		decrByCmd.Execute()
	}
}

