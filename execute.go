package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shubhdevelop/YAKVS/command"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/snapshot"
	"github.com/shubhdevelop/YAKVS/store"
)

func ExecuteCommand(cmd *parser.Command, store *store.Store) {
	fmt.Println("Executing command:", cmd)
	switch strings.ToUpper(cmd.Name) {
	case "BGSAVE":
		if len(cmd.Args) > 0 {
			fmt.Println("Error: BGSAVE doesn't require any arguments")
		}

		fmt.Println("snapshot started")
		snapshot.Start()
		fmt.Println("+OK\r")
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
		if len(cmd.Args) < 1 {
			fmt.Println("Error: EXISTS requires 1 argument (key)")
			return
		}
		exists := store.Exists(cmd.Args[0])
		if exists {
			fmt.Println(":1\r")
		} else {
			fmt.Println(":0\r")
		}
	case "TTL":
		if len(cmd.Args) < 1 {
			fmt.Println("Error: TTL requires 1 argument (key)")
			return
		}
		ttl := store.GetTTL(cmd.Args[0])
		fmt.Printf(":%d\r\n", ttl)
	case "EXPIRE":
		if len(cmd.Args) < 2 {
			fmt.Println("Error: EXPIRE requires 2 arguments (key, ttl)")
			return
		}
		ttl, err := strconv.ParseInt(cmd.Args[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing TTL:", err)
			return
		}
		// calculate the unix timestamp for the given ttl
		ttl = time.Now().Unix() + ttl
		success := store.SetTTL(cmd.Args[0], ttl)
		if success {
			fmt.Println("+OK\r")
		} else {
			fmt.Println(":0\r")
		}
	case "EXPIREAT":
		if len(cmd.Args) < 2 {
			fmt.Println("Error: EXPIREAT requires 2 arguments (key, timestamp)")
			return
		}
		// we expect the ttl to be in unix timestamp
		ttl, err := strconv.ParseInt(cmd.Args[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing TTL:", err)
			return
		}
		success := store.SetTTL(cmd.Args[0], ttl)
		if success {
			fmt.Println("+OK\r")
		} else {
			fmt.Println(":0\r")
		}
	case "PERSIST":
		if len(cmd.Args) < 1 {
			fmt.Println("Error: PERSIST requires 1 argument (key)")
			return
		}
		success := store.RemoveExpiry(cmd.Args[0])
		if success {
			fmt.Println(":1\r")
		} else {
			fmt.Println(":0\r")
		}
	}
}

