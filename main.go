package main

import (
	"bufio"
	"fmt"
	"github.com/shubhdevelop/YAKVS/scanner"
	"log"
	"os"
	"strings"
)

var File *os.File

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

		// Process the input
		processedLine := preprocessInput(line[:len(line)-1])

		// Try to detect if it's RESP format
		if isRESPFormat(processedLine) {
			// Create a RESPReader for this specific input
			respReader := NewRESPReader(strings.NewReader(processedLine))
			command, err := respReader.ReadCommand()
			if err != nil {
				fmt.Printf("Error parsing RESP command: %v\n", err)
			} else {
				fmt.Printf("RESP Command received:\n%s\n", command)
				// Also process with scanner for additional analysis
				fmt.Println("Scanner analysis:")
				_, err := File.WriteString(line)
				if err != nil {
					log.Fatalf("failed to write to file: %v", err)
				}
				processWithScanner(command)
			}
		} else {
			// Process as regular input using scanner
			processWithScanner(processedLine)
		}
	}
}

func processWithScanner(input string) {
	scanner := scanner.Scanner{
		Source: input,
	}
	tokens, err := scanner.ScanTokens()
	if err != nil {
		fmt.Println("Error scanning tokens:", err)
		return
	}
	fmt.Printf("Tokens (%d):\n", len(tokens))
	for i, token := range tokens {
		fmt.Printf("  %d: %v\n", i, token)
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
	file, err := os.OpenFile("base.aof", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	File = file
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
}

func main() {
	fmt.Println("YAKVS")
	runPrompt()
	defer File.Close()
}
