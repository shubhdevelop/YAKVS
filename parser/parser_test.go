package parser

import (
	"testing"
)

func TestNewStreamingParser(t *testing.T) {
	data := []byte("test data")
	parser := NewStreamingParser(data)
	
	if parser == nil {
		t.Fatal("Expected parser to be created, got nil")
	}
	
	if parser.pos != 0 {
		t.Errorf("Expected pos to be 0, got %d", parser.pos)
	}
	
	if parser.len != len(data) {
		t.Errorf("Expected len to be %d, got %d", len(data), parser.len)
	}
	
	if string(parser.buf) != string(data) {
		t.Errorf("Expected buf to be %s, got %s", string(data), string(parser.buf))
	}
}

func TestReadUntilCRLF(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		hasError bool
	}{
		{
			name:     "valid CRLF",
			input:    []byte("hello\r\nworld"),
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "no CRLF",
			input:    []byte("hello world"),
			expected: nil,
			hasError: true,
		},
		{
			name:     "only CR",
			input:    []byte("hello\rworld"),
			expected: nil,
			hasError: true,
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: nil,
			hasError: true,
		},
		{
			name:     "only CRLF",
			input:    []byte("\r\n"),
			expected: []byte(""),
			hasError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result, err := parser.readUntilCRLF()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if string(result) != string(tt.expected) {
					t.Errorf("Expected %s, got %s", string(tt.expected), string(result))
				}
			}
		})
	}
}

func TestSkipCRLF(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		hasError bool
	}{
		{
			name:     "valid CRLF",
			input:    []byte("\r\n"),
			hasError: false,
		},
		{
			name:     "no CRLF",
			input:    []byte("hello"),
			hasError: true,
		},
		{
			name:     "only CR",
			input:    []byte("\r"),
			hasError: true,
		},
		{
			name:     "empty input",
			input:    []byte(""),
			hasError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			err := parser.skipCRLF()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestReadBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		hasError bool
	}{
		{
			name:     "valid bulk string",
			input:    []byte("$5\r\nhello\r\n"),
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "empty bulk string",
			input:    []byte("$0\r\n\r\n"),
			expected: []byte(""),
			hasError: false,
		},
		{
			name:     "invalid length",
			input:    []byte("$abc\r\nhello\r\n"),
			expected: nil,
			hasError: true,
		},
		{
			name:     "missing CRLF after length",
			input:    []byte("$5hello\r\n"),
			expected: nil,
			hasError: true,
		},
		{
			name:     "missing CRLF after content",
			input:    []byte("$5\r\nhello"),
			expected: nil,
			hasError: true,
		},
		{
			name:     "content exceeds buffer",
			input:    []byte("$10\r\nhello"),
			expected: nil,
			hasError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result, err := parser.readBulkString()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if string(result) != string(tt.expected) {
					t.Errorf("Expected %s, got %s", string(tt.expected), string(result))
				}
			}
		})
	}
}

func TestReadSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		hasError bool
	}{
		{
			name:     "valid simple string",
			input:    []byte("hello\r\n"),
			expected: []byte("hello"),
			hasError: false,
		},
		{
			name:     "empty simple string",
			input:    []byte("\r\n"),
			expected: []byte(""),
			hasError: false,
		},
		{
			name:     "no CRLF",
			input:    []byte("hello"),
			expected: nil,
			hasError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result, err := parser.readSimpleString()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if string(result) != string(tt.expected) {
					t.Errorf("Expected %s, got %s", string(tt.expected), string(result))
				}
			}
		})
	}
}

func TestParseBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *Command
	}{
		{
			name:  "valid bulk string",
			input: []byte("$5\r\nhello\r\n"),
			expected: &Command{
				Name: "hello",
				Args: []string{},
			},
		},
		{
			name:  "empty bulk string",
			input: []byte("$0\r\n\r\n"),
			expected: &Command{
				Name: "",
				Args: []string{},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result := parser.parseBulkString()
			
			if result.Name != tt.expected.Name {
				t.Errorf("Expected name %s, got %s", tt.expected.Name, result.Name)
			}
			
			if len(result.Args) != len(tt.expected.Args) {
				t.Errorf("Expected %d args, got %d", len(tt.expected.Args), len(result.Args))
			}
		})
	}
}

func TestParseSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *Command
	}{
		{
			name:  "valid simple string",
			input: []byte("+hello\r\n"),
			expected: &Command{
				Name: "hello",
				Args: []string{},
			},
		},
		{
			name:  "empty simple string",
			input: []byte("+\r\n"),
			expected: &Command{
				Name: "",
				Args: []string{},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result := parser.parseSimpleString()
			
			if result.Name != tt.expected.Name {
				t.Errorf("Expected name %s, got %s", tt.expected.Name, result.Name)
			}
			
			if len(result.Args) != len(tt.expected.Args) {
				t.Errorf("Expected %d args, got %d", len(tt.expected.Args), len(result.Args))
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *Command
	}{
		{
			name:  "array with bulk strings",
			input: []byte("*2\r\n$4\r\nPING\r\n$4\r\ntest\r\n"),
			expected: &Command{
				Name: "PING",
				Args: []string{"test"},
			},
		},
		{
			name:  "array with simple strings",
			input: []byte("*2\r\n+OK\r\n+SUCCESS\r\n"),
			expected: &Command{
				Name: "OK",
				Args: []string{"SUCCESS"},
			},
		},
		{
			name:  "empty array",
			input: []byte("*0\r\n"),
			expected: &Command{
				Name: "",
				Args: []string{},
			},
		},
		{
			name:  "single element array",
			input: []byte("*1\r\n$4\r\nPING\r\n"),
			expected: &Command{
				Name: "PING",
				Args: []string{},
			},
		},
		{
			name:  "SET command with multiple args",
			input: []byte("*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nEX\r\n$3\r\n300\r\n"),
			expected: &Command{
				Name: "SET",
				Args: []string{"key", "value", "EX", "300"},
			},
		},
		{
			name:  "command with nested array as argument",
			input: []byte("*3\r\n$4\r\nEVAL\r\n$3\r\nlua\r\n*2\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"),
			expected: &Command{
				Name: "EVAL",
				Args: []string{"lua", "key", "value"},
			},
		},
		{
			name:  "GEOADD command with multiple args",
			input: []byte("*6\r\n$6\r\nGEOADD\r\n$6\r\nplaces\r\n$2\r\nXX\r\n$9\r\n139.69171\r\n$10\r\n35.6895000\r\n$5\r\nTokyo\r\n"),
			expected: &Command{
				Name: "GEOADD",
				Args: []string{"places", "XX", "139.69171", "35.6895000", "Tokyo"},
			},
		},
		
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result := parser.ParseArray()
			
			if result.Name != tt.expected.Name {
				t.Errorf("Expected name %s, got %s", tt.expected.Name, result.Name)
			}
			
			if len(result.Args) != len(tt.expected.Args) {
				t.Errorf("Expected %d args, got %d", len(tt.expected.Args), len(result.Args))
			}
			
			for i, arg := range result.Args {
				if i < len(tt.expected.Args) && arg != tt.expected.Args[i] {
					t.Errorf("Expected arg[%d] %s, got %s", i, tt.expected.Args[i], arg)
				}
			}
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *Command
		hasError bool
	}{
		{
			name:  "array command",
			input: []byte("*2\r\n$4\r\nPING\r\n$4\r\ntest\r\n"),
			expected: &Command{
				Name: "PING",
				Args: []string{"test"},
			},
			hasError: false,
		},
		{
			name:  "bulk string command",
			input: []byte("$4\r\nPING\r\n"),
			expected: &Command{
				Name: "PING",
				Args: []string{},
			},
			hasError: false,
		},
		{
			name:     "unexpected token",
			input:    []byte("invalid"),
			expected: nil,
			hasError: true,
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: nil,
			hasError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewStreamingParser(tt.input)
			result, err := parser.ParseCommand()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				
				if result == nil {
					t.Errorf("Expected result to be non-nil")
					return
				}
				
				if result.Name != tt.expected.Name {
					t.Errorf("Expected name %s, got %s", tt.expected.Name, result.Name)
				}
				
				if len(result.Args) != len(tt.expected.Args) {
					t.Errorf("Expected %d args, got %d", len(tt.expected.Args), len(result.Args))
				}
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("very long bulk string", func(t *testing.T) {
		longString := make([]byte, 1000)
		for i := range longString {
			longString[i] = 'a'
		}
		
		input := append([]byte("$1000\r\n"), longString...)
		input = append(input, []byte("\r\n")...)
		
		parser := NewStreamingParser(input)
		result, err := parser.readBulkString()
		
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if len(result) != 1000 {
			t.Errorf("Expected length 1000, got %d", len(result))
		}
	})
	
	t.Run("array with mixed types", func(t *testing.T) {
		input := []byte("*3\r\n$4\r\nPING\r\n+OK\r\n$5\r\nworld\r\n")
		parser := NewStreamingParser(input)
		result := parser.ParseArray()
		
		expected := &Command{
			Name: "PING",
			Args: []string{"OK", "world"},
		}
		
		if result.Name != expected.Name {
			t.Errorf("Expected name %s, got %s", expected.Name, result.Name)
		}
		
		if len(result.Args) != len(expected.Args) {
			t.Errorf("Expected %d args, got %d", len(expected.Args), len(result.Args))
		}
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("incomplete bulk string", func(t *testing.T) {
		input := []byte("$5\r\nhel")
		parser := NewStreamingParser(input)
		_, err := parser.readBulkString()
		
		if err == nil {
			t.Errorf("Expected error for incomplete bulk string")
		}
	})
	
	t.Run("missing CRLF in simple string", func(t *testing.T) {
		input := []byte("hello")
		parser := NewStreamingParser(input)
		_, err := parser.readSimpleString()
		
		if err == nil {
			t.Errorf("Expected error for missing CRLF")
		}
	})
	
	t.Run("invalid bulk string length", func(t *testing.T) {
		input := []byte("$abc\r\nhello\r\n")
		parser := NewStreamingParser(input)
		_, err := parser.readBulkString()
		
		if err == nil {
			t.Errorf("Expected error for invalid bulk string length")
		}
	})
}
