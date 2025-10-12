package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/shubhdevelop/YAKVS/aof"
	"github.com/shubhdevelop/YAKVS/executor"
	"github.com/shubhdevelop/YAKVS/parser"
	"github.com/shubhdevelop/YAKVS/store"
	"github.com/shubhdevelop/YAKVS/utils"
)

func StartServer(port string, kvStore *store.Store, aofManager *aof.AOFManager) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server started on port %s\n", port)	

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v", err)
		}
		go handleConnection(conn, kvStore, aofManager)
	}
}

func handleConnection(conn net.Conn, kvStore *store.Store, aofManager *aof.AOFManager) {
	defer conn.Close()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n') // resp read the line
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}
		if line == "\n" || line == "" {
			continue
		}
		fmt.Println("Received message:", line)
		resp, err := utils.ToRESP(line)
		if err != nil {
			fmt.Printf("Error converting to RESP: %v\n", err)
			continue
		}

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
		resp, err := executor.ExecuteCommand(command, kvStore)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			continue
		}
		fmt.Println("command response:", resp)
		// write the response to the client
		conn.Write([]byte(resp))
	}
	}
}

