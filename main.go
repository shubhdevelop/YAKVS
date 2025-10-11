package main

import (
	"fmt"
	"log"

	"github.com/shubhdevelop/YAKVS/aof"
	"github.com/shubhdevelop/YAKVS/executor"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/server"
	"github.com/shubhdevelop/YAKVS/store"
)

var aofManager *aof.AOFManager
var KvStore *store.Store


func init() {
	// Initialize AOF manager
	aofManager = aof.NewAOFManager("base.aof")
	err := aofManager.Initialize()
	if err != nil {
		log.Fatalf("Error initializing AOF manager: %v", err)
	}
	// Initialize store
	KvStore = store.NewStore()
}

func main() {
	fmt.Println("YAKVS")
	// Read and execute commands from AOF file
	err := aofManager.ReadAndExecuteCommands(func(cmd *parser.Command) {
		executor.ExecuteCommand(cmd, KvStore)
	})
	
	if err != nil {
		log.Fatalf("Error reading AOF file: %v", err)
	}

	server.StartServer("8080", KvStore, aofManager)

	// runPrompt()
	defer aofManager.Close()
}
