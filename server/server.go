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

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n') // resp read the line
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}
		if line == "\n" || line == "" {
			fmt.Println("Empty line received")
			conn.Write([]byte("$-1\r\n\n"))
			continue
		}
		resp, err := utils.ToRESP(line)
		if err != nil {
			fmt.Printf("Error converting to RESP: %v\n", err)
			conn.Write([]byte("$-1\r\n\n"))
			continue
		}

		if utils.IsRESPFormat(resp) {
			parser := parser.NewStreamingParser([]byte(resp))
			command, err := parser.ParseCommand()
			if err != nil {
				fmt.Printf("Error parsing RESP command: %v\n", err)
				conn.Write([]byte("$-1\r\n\n"))
				continue
			}


			if aofManager.ShouldPersistCommand(command.Name) {
				err := aofManager.WriteCommand(resp)
				if err != nil {
					fmt.Printf("failed to write to AOF file: %v\n", err)
					conn.Write([]byte("$-1\r\n\n"))
					continue
				}
				fmt.Println("Writing command to AOF file:", resp)
			}

			resultChan := make(chan executor.ResultWithError, 1) 
			executor.ExecuteCommandAysnc(command, kvStore, resultChan)
			
			result := <-resultChan
			if result.Err != nil {
				fmt.Printf("error executing command: %v\n", result.Err)
				conn.Write([]byte("$-1\r\n\n"))
				continue
			}
			fmt.Printf("command response: %s\n", result.Result)
			conn.Write([]byte(result.Result + "\n"))
		}
	}
}
