package main

import (
	"fmt"

	"github.com/shubhdevelop/YAKVS/parser"
)

func ExecuteCommand(command *parser.Command) {
	fmt.Println("Executing command:", command)
	switch command.Name {
	case "SET":
		kvStore.SetValue(command.Args[0], command.Args[1])
	case "GET":
		fmt.Println(kvStore.GetValue(command.Args[0]))
	case "DEL":
		kvStore.DeleteValue(command.Args[0])
	case "EXISTS":
		fmt.Println(kvStore.Exists(command.Args[0]))
	case "TTL":
		fmt.Println(kvStore.GetTTL(command.Args[0]))
	}
}