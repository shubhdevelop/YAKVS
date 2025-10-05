package main

import (
	"bufio"
	"fmt"
	"github.com/shubhdevelop/YAKVS/scanner"
	"os"
	"strings"
)

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error Reading the line")
			continue
		} else if line == "\n" || line == "" {
			continue
		} else if line == "clear\n" {
			fmt.Print("\033[H\033[2J")
			continue
		} else if line == "exit\n" {
			break
		}

		// Convert literal escape sequences to actual control characters for interactive input
		processedLine := preprocessInput(line[:len(line)-2])
		
		scanner := scanner.Scanner{
			Source: processedLine,
		}
		tokens, err := scanner.ScanTokens()
		if err != nil {
			fmt.Println("Error scanning tokens:", err)
			continue
		}
		fmt.Printf("Tokens (%d):\n", len(tokens))
		for i, token := range tokens {
			fmt.Printf("  %d: %v\n", i, token)
		}
	}
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

func main() {
	fmt.Println("YAKVS")
	runPrompt()
}
