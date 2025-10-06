package main

import (
	"fmt"
	"github.com/shubhdevelop/YAKVS/parser"
)

func ExecuteCommand(command *parser.Command) {
	fmt.Println("Executing command:", command)
	switch command.Name {
	case "SET":
		store.SetValue(command.Args[0], command.Args[1])
	case "GET":
		fmt.Println(store.GetValue(command.Args[0]))
	case "DEL":
		store.DeleteValue(command.Args[0])
	case "EXISTS":
		fmt.Println(store.Exists(command.Args[0]))
	}
}