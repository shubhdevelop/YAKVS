package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
)

func ExecuteCommand(command *parser.Command, store *store.Store) {
	fmt.Println("Executing command:", command)
	switch strings.ToUpper(command.Name) {
	case "SET":
		if len(command.Args) < 2 {
			fmt.Println("Error: SET requires 2 arguments (key, value)")
			return
		}
		store.SetValue(command.Args[0], command.Args[1])
		fmt.Println("+OK\r")
	case "GET":
		if len(command.Args) < 1 {
			fmt.Println("Error: GET requires 1 argument (key)")
			return
		}
		value := store.GetValue(command.Args[0])
		if value == nil {
			fmt.Println("$-1\r")
		} else {
			fmt.Printf("$%d\r\n%s\r\n", len(fmt.Sprintf("%v", value)), value)
		}
	case "DEL":
		if len(command.Args) < 1 {
			fmt.Println("Error: DEL requires 1 argument (key)")
			return
		}
		success := store.DeleteValue(command.Args[0])
		if success {
			fmt.Println("+OK\r")
		} else {
			fmt.Println("$-1\r")
		}
	case "EXISTS":
		if len(command.Args) < 1 {
			fmt.Println("Error: EXISTS requires 1 argument (key)")
			return
		}
		exists := store.Exists(command.Args[0])
		if exists {
			fmt.Println(":1\r")
		} else {
			fmt.Println(":0\r")
		}
	case "TTL":
		if len(command.Args) < 1 {
			fmt.Println("Error: TTL requires 1 argument (key)")
			return
		}
		ttl := store.GetTTL(command.Args[0])
		fmt.Printf(":%d\r\n", ttl)
	case "EXPIRE":
		if len(command.Args) < 2 {
			fmt.Println("Error: EXPIRE requires 2 arguments (key, ttl)")
			return
		}
		ttl, err := strconv.ParseInt(command.Args[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing TTL:", err)
			return
		}
		// calculate the unix timestamp for the given ttl
		ttl = time.Now().Unix() + ttl
		success := store.SetTTL(command.Args[0], ttl)
		if success {
			fmt.Println("+OK\r")
		} else {
			fmt.Println(":0\r")
		}
	case "EXPIREAT":
		if len(command.Args) < 2 {
			fmt.Println("Error: EXPIREAT requires 2 arguments (key, timestamp)")
			return
		}
		// we expect the ttl to be in unix timestamp
		ttl, err := strconv.ParseInt(command.Args[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing TTL:", err)
			return
		}
		success := store.SetTTL(command.Args[0], ttl)
		if success {
			fmt.Println("+OK\r")
		} else {
			fmt.Println(":0\r")
		}
	}
}