package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func runPrompt(conn net.Conn) {
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

		fmt.Println("Sending message:", line)
		// resp, err := utils.ToRESP(line[:len(line)-1])
		resp := line

		if resp != "" {
			fmt.Println("Writing message:", resp)
			conn.Write([]byte(resp))
		}
		reader := bufio.NewReader(conn)
		line, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			continue
		}
		fmt.Println("Received message from server:", line)
		if line == "exit\n" {
			break
		}
	}
}


func Connect(address string) (*net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}


func main() {
	conn, err := Connect("localhost:8080")
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer (*conn).Close()
	
	// Start the interactive prompt
	runPrompt(*conn)
}


// it's an cli client