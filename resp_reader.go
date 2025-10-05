package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RESPReader handles reading complete RESP commands
type RESPReader struct {
	reader *bufio.Reader
}

// NewRESPReader creates a new RESP reader
func NewRESPReader(reader io.Reader) *RESPReader {
	return &RESPReader{
		reader: bufio.NewReader(reader),
	}
}

// ReadCommand reads a complete RESP command
func (r *RESPReader) ReadCommand() (string, error) {
	var result strings.Builder
	
	// Read the first character to determine the type
	firstChar, err := r.reader.ReadByte()
	if err != nil {
		return "", err
	}
	
	result.WriteByte(firstChar)
	
	switch firstChar {
	case '*': // Array
		return r.readArray(&result)
	case '$': // Bulk string
		return r.readBulkString(&result)
	case '+': // Simple string
		return r.readSimpleString(&result)
	case '-': // Error
		return r.readSimpleString(&result)
	case ':': // Integer
		return r.readSimpleString(&result)
	case '!': // BLOB
		return r.readBulkString(&result)
	case '=': // Verbatim string
		return r.readBulkString(&result)
	case '%': // Map
		return r.readMap(&result)
	case '~': // Set
		return r.readSet(&result)
	case '>': // Push data
		return r.readPushData(&result)
	case '_': // Null
		return r.readSimpleString(&result)
	default:
		return "", fmt.Errorf("unknown RESP type: %c", firstChar)
	}
}

// readArray reads a RESP array (also used for Map, Set, Push data)
func (r *RESPReader) readArray(result *strings.Builder) (string, error) {
	// Read the array length
	length, err := r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString(length)
	
	// Parse the length
	arrayLen, err := strconv.Atoi(length)
	if err != nil {
		return "", err
	}
	
	// Read each element
	for i := 0; i < arrayLen; i++ {
		// Read the element type
		elementType, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		result.WriteByte(elementType)
		
		// Handle different element types
		switch elementType {
		case '$': // Bulk string
			element, err := r.readBulkString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '!': // BLOB
			element, err := r.readBulkString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '=': // Verbatim string
			element, err := r.readBulkString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '*': // Nested array
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '%': // Nested map
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '~': // Nested set
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '>': // Nested push data
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		default: // Simple string, error, integer, null
			element, err := r.readSimpleString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		}
	}
	
	return result.String(), nil
}

// readMap reads a RESP map (key-value pairs)
func (r *RESPReader) readMap(result *strings.Builder) (string, error) {
	// Read the map length (number of key-value pairs)
	length, err := r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString(length)
	
	// Parse the length
	mapLen, err := strconv.Atoi(length)
	if err != nil {
		return "", err
	}
	
	// Read key-value pairs (2 * mapLen elements)
	for i := 0; i < mapLen*2; i++ {
		// Read the element type
		elementType, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		result.WriteByte(elementType)
		
		// Handle different element types
		switch elementType {
		case '$', '!', '=': // Bulk string, BLOB, Verbatim string
			element, err := r.readBulkString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '*', '%', '~', '>': // Nested structures
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		default: // Simple string, error, integer, null
			element, err := r.readSimpleString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		}
	}
	
	return result.String(), nil
}

// readSet reads a RESP set (unique elements)
func (r *RESPReader) readSet(result *strings.Builder) (string, error) {
	// Read the set length
	length, err := r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString(length)
	
	// Parse the length
	setLen, err := strconv.Atoi(length)
	if err != nil {
		return "", err
	}
	
	// Read set elements
	for i := 0; i < setLen; i++ {
		// Read the element type
		elementType, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		result.WriteByte(elementType)
		
		// Handle different element types
		switch elementType {
		case '$', '!', '=': // Bulk string, BLOB, Verbatim string
			element, err := r.readBulkString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '*', '%', '~', '>': // Nested structures
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		default: // Simple string, error, integer, null
			element, err := r.readSimpleString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		}
	}
	
	return result.String(), nil
}

// readPushData reads a RESP push data
func (r *RESPReader) readPushData(result *strings.Builder) (string, error) {
	// Read the push data length
	length, err := r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString(length)
	
	// Parse the length
	pushLen, err := strconv.Atoi(length)
	if err != nil {
		return "", err
	}
	
	// Read push data elements
	for i := 0; i < pushLen; i++ {
		// Read the element type
		elementType, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		result.WriteByte(elementType)
		
		// Handle different element types
		switch elementType {
		case '$', '!', '=': // Bulk string, BLOB, Verbatim string
			element, err := r.readBulkString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		case '*', '%', '~', '>': // Nested structures
			element, err := r.readArray(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		default: // Simple string, error, integer, null
			element, err := r.readSimpleString(&strings.Builder{})
			if err != nil {
				return "", err
			}
			result.WriteString(element)
		}
	}
	
	return result.String(), nil
}

// readBulkString reads a RESP bulk string
func (r *RESPReader) readBulkString(result *strings.Builder) (string, error) {
	// Read the length
	length, err := r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString(length)
	
	// Parse the length
	stringLen, err := strconv.Atoi(length)
	if err != nil {
		return "", err
	}
	
	// Read the string content
	content := make([]byte, stringLen)
	_, err = io.ReadFull(r.reader, content)
	if err != nil {
		return "", err
	}
	result.Write(content)
	
	// Read the final CRLF
	_, err = r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString("\r\n")
	
	return result.String(), nil
}

// readSimpleString reads a RESP simple string
func (r *RESPReader) readSimpleString(result *strings.Builder) (string, error) {
	content, err := r.readUntilCRLF()
	if err != nil {
		return "", err
	}
	result.WriteString(content)
	
	return result.String(), nil
}

// readUntilCRLF reads until \r\n
func (r *RESPReader) readUntilCRLF() (string, error) {
	var result strings.Builder
	
	for {
		char, err := r.reader.ReadByte()
		if err != nil {
			return "", err
		}
		
		if char == '\r' {
			// Check for \n
			nextChar, err := r.reader.ReadByte()
			if err != nil {
				return "", err
			}
			if nextChar == '\n' {
				break
			}
			// Put back the character if it's not \n
			r.reader.UnreadByte()
		}
		
		result.WriteByte(char)
	}
	
	return result.String(), nil
}
