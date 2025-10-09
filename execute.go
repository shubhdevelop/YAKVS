package main

import (
	"fmt"
	"strings"

	"github.com/shubhdevelop/YAKVS/command"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

func ExecuteCommand(cmd *parser.Command, store *store.Store) {
	fmt.Println("Executing command:", cmd)
	switch strings.ToUpper(cmd.Name) {
	case "BGSAVE":
		bgSaveCmd := command.NewBgSaveCommand(cmd, store)
		bgSaveCmd.Execute()
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
			return
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
	}
}

