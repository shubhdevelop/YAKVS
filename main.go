package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"github.com/shubhdevelop/YAKVS/parser"
)

var WriteFile *os.File
var ReadFile *os.File
var store *Store

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

		// Try to detect if it's RESP format
		if isRESPFormat(line[:len(line)-1]) {
			// Preprocess input to convert literal \r\n to actual control characters
			processedInput := preprocessInput(line[:len(line)-1])
			parser := parser.NewStreamingParser([]byte(processedInput))
			fmt.Println("Parsing RESP command:", line[:len(line)-1])
			command, err := parser.ParseCommand()
			if err != nil {
				fmt.Printf("Error parsing RESP command: %v\n", err)
			}
			if command.Name == "SET" {
				_, err := WriteFile.WriteString(processedInput)

				if err != nil {
					log.Fatalf("failed to write to file: %v", err)
				}
				ExecuteCommand(command)
			}
			fmt.Println("Executed command:", command)
		}
	}
}


func isRESPFormat(input string) bool {
	// Check if input starts with RESP protocol indicators
	if len(input) == 0 {
		return false
	}

	firstChar := input[0]
	// RESP protocol starts with specific characters
	respTypes := []byte{'*', '$', '+', '-', ':', '!', '=', '%', '~', '>', '_'}

	for _, respType := range respTypes {
		if firstChar == respType {
			return true
		}
	}

	return false
}

func preprocessInput(input string) string {
	// Convert literal escape sequences to actual control characters
	// This is only needed for interactive input where user types \r\n literally
	result := input
	result = strings.ReplaceAll(result, "\\r", "\r")
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\t", "\t")
	result = strings.ReplaceAll(result, "\\\\", "\\")
	return result
}

func init() {
	// Open file for writing (AOF - Append Only File)
	writeFile, err := os.OpenFile("base.aof", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening write file:", err)
		return
	}
	WriteFile = writeFile

	// Open file for reading
	readFile, err := os.Open("base.aof")
	if err != nil {
		// If file doesn't exist yet, that's okay - we'll create it when we write
		if !os.IsNotExist(err) {
			fmt.Println("Error opening read file:", err)
		}
		ReadFile = nil
		return
	}
	ReadFile = readFile
	store = &Store{
		Values: make(map[string]interface{}),
	}
}

func main() {
	fmt.Println("YAKVS")
	if ReadFile != nil {
		fmt.Println("Reading from AOF file:")
		// Read the entire file content
		fileContent, err := io.ReadAll(ReadFile)
		if err != nil {
			log.Fatalf("Error reading AOF file: %v", err)
		}
		
		// Parse the entire file as RESP commands
		parser := parser.NewStreamingParser(fileContent)
		
		// Parse all commands in the file
		for {
			command, err := parser.ParseCommand()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error parsing RESP command: %v\n", err)
				break
			}
			ExecuteCommand(command)
		}
		ReadFile.Close()
	} else {
		fmt.Println("No AOF file found, starting fresh.")
	}

	runPrompt()
	defer WriteFile.Close()
}
