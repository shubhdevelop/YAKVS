package utils

import "strings"

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


// IsRESPFormat checks if input starts with RESP protocol indicators
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
