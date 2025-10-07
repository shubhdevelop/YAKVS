package utils

import (
	"fmt"
	"strings"
)

// Deprecated: No longer needed because we turning interactive input to RESP protocol
// PreprocessInput converts literal escape sequences to actual control characters
// This is only needed for interactive input where user types \r\n literally
func PreprocessInput(input string) string {
	result := input
	result = strings.ReplaceAll(result, "\\r", "\r")
	result = strings.ReplaceAll(result, "\\n", "\n")
	result = strings.ReplaceAll(result, "\\t", "\t")
	result = strings.ReplaceAll(result, "\\\\", "\\")
	return result
}

// Checks if input starts with RESP protocol indicators
func IsRESPFormat(input string) bool {
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

func ToRESP(command string) (string, error) {
	parts := strings.Fields(command) // Split command by spaces

	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := strings.ToUpper(parts[0]) // Get the command name (e.g., "SET")

	var respBuilder strings.Builder

	switch cmd {
	case "EXPIRE", "EXPIREAT":
		if len(parts) < 3 {
			return "", fmt.Errorf("%s command requires at least two arguments: key and TTL", cmd)
		}
		// The number of arguments in the array will be 1 (command) + key + TTL
		respBuilder.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))
		for _, part := range parts {
			respBuilder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(part), part))
		}
		return respBuilder.String(), nil
	case "SET" :
		// SET key value [EX seconds] [PX milliseconds] [NX|XX]
		if len(parts) < 3 {
			return "", fmt.Errorf("SET command requires at least a key and a value")
		}
		// The number of arguments in the array will be 1 (SET) + key + value + options
		respBuilder.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))
		for _, part := range parts {
			respBuilder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(part), part))
		}
		return respBuilder.String(), nil
	
	// all uppercase commands or all lowercase commands both are valid
	case "GET", "DEL", "EXISTS", "TTL" : 
		if len(parts) < 2 {
			return "", fmt.Errorf("%s command requires at least one key", cmd)
		}
		// For GET, DEL, EXISTS, TTL, the number of arguments is 1 (command) + number of keys
		respBuilder.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))
		for _, part := range parts {
			respBuilder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(part), part))
		}
		return respBuilder.String(), nil

	default:
		return "", fmt.Errorf("unsupported command: %s", cmd)
	}
}
