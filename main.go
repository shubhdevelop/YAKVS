package main

import (
	"bufio"
	"fmt"
	"os"
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
	}
}

func main() {

	fmt.Println("YAKVS")
	runPrompt()
}
