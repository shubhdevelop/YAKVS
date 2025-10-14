package executor

import (
	"fmt"
	"strings"

	"github.com/shubhdevelop/YAKVS/command"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

type ResultWithError struct {
	Result string
	Err    error
}

func ExecuteCommandAysnc(cmd *parser.Command, store *store.Store, ch chan ResultWithError) {
	fmt.Println("Executing command:", cmd)
	
	// Execute command concurrently in a goroutine
	go func() {
		store.Mu.Lock()
		defer store.Mu.Unlock()
		
		switch strings.ToUpper(cmd.Name) {
		case "BGSAVE":
			bgSaveCmd := command.NewBgSaveCommand(cmd, store)
			ch <- ResultWithError{Result: bgSaveCmd.Execute(), Err: nil}
		case "SET":
			setCmd := command.NewSetCommand(cmd, store)
			ch <- ResultWithError{Result: setCmd.Execute(), Err: nil}
		case "GET":
			getCmd := command.NewGetCommand(cmd, store)
			ch <- ResultWithError{Result: getCmd.Execute(), Err: nil}
		case "DEL":
			delCmd := command.NewDelCommand(cmd, store)
			ch <- ResultWithError{Result: delCmd.Execute(), Err: nil}
		case "EXISTS":
			existsCmd := command.NewExistsCommand(cmd, store)
			ch <- ResultWithError{Result: existsCmd.Execute(), Err: nil}
		case "TTL":
			ttlCmd := command.NewTtlCommand(cmd, store)
			ch <- ResultWithError{Result: ttlCmd.Execute(), Err: nil}
		case "EXPIRE":
			expireCmd := command.NewExpireCommand(cmd, store)
			ch <- ResultWithError{Result: expireCmd.Execute(), Err: nil}
		case "EXPIREAT":
			expireAtCmd := command.NewExpireAtCommand(cmd, store)
			ch <- ResultWithError{Result: expireAtCmd.Execute(), Err: nil}
		case "PERSIST":
			persistCmd := command.NewPersistCommand(cmd, store)
			ch <- ResultWithError{Result: persistCmd.Execute(), Err: nil}
		case "INCRBY":
			incrByCmd := command.NewIncreByCommand(cmd, store)
			ch <- ResultWithError{Result: incrByCmd.Execute(), Err: nil}
		case "DECRBY":
			decrByCmd := command.NewDecreByCommand(cmd, store)
			ch <- ResultWithError{Result: decrByCmd.Execute(), Err: nil}
		case "INCR":
			incrByCmd := command.NewIncrCommand(cmd, store)
			ch <- ResultWithError{Result: incrByCmd.Execute(), Err: nil}
		case "DECR":
			decrByCmd := command.NewDecrCommand(cmd, store)
			ch <- ResultWithError{Result: decrByCmd.Execute(), Err: nil}
		default:
			ch <- ResultWithError{Result: "", Err: fmt.Errorf("invalid command: %s", cmd.Name)}
		}
	}()
}

func ExecuteCommandSync(cmd *parser.Command, store *store.Store) {
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
