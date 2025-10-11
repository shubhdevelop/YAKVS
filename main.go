package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/shubhdevelop/YAKVS/aof"
	"github.com/shubhdevelop/YAKVS/executor"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/server"
	"github.com/shubhdevelop/YAKVS/store"
	"github.com/shubhdevelop/YAKVS/utils"
)

var aofManager *aof.AOFManager
var KvStore *store.Store


// Todo: build a client for the yakvs server
func runPrompt() {
	// Use regular reader for line-by-line input
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")

		// Read line by line first
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("\nGoodbye!")
				break
			}
			fmt.Println("Error Reading the line:", err)
			continue
		} else if line == "\n" || line == "" {
			continue
		} else if line == "clear\n" {
			fmt.Print("\033[H\033[2J")
			continue
		} else if line == "exit\n" {
			break
		}

		resp, err := utils.ToRESP(line[:len(line)-1])
		if err != nil {
			fmt.Printf("Error converting to RESP: %v\n", err)
			continue
		}
		// Try to detect if it's RESP format
		if utils.IsRESPFormat(resp) {
				// Preprocess input to convert literal \r\n to actual control characters
			// processedInput := utils.PreprocessInput(resp)
			parser := parser.NewStreamingParser([]byte(resp))
			fmt.Println("Parsing RESP command:", resp)
			command, err := parser.ParseCommand()
			if err != nil {
				fmt.Printf("Error parsing RESP command: %v\n", err)
			}
			// check if command should be persisted
			if aofManager.ShouldPersistCommand(command.Name) {
				err := aofManager.WriteCommand(resp)
				if err != nil {
					log.Fatalf("failed to write to AOF file: %v", err)
				}
			}
			executor.ExecuteCommand(command, KvStore)
		}
	}
}

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
