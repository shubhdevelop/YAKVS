package parser

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type StreamingParser struct {
	buf []byte
	pos int
	len int
}

type Command struct {
	Name string
	Args []string
}

func (p *StreamingParser) ParseCommand() (*Command, error) {
	if p.pos >= p.len {
		return nil, io.EOF
	}

	switch p.buf[p.pos] {
		case '*', '~', '>', '%': 
			return p.ParseArray(), nil
		case '$':
			return p.parseBulkString(), nil
		default:
			return nil, fmt.Errorf("unexpected token: %c", p.buf[p.pos])
	}
}


func (p *StreamingParser) ParseArray() *Command {
	p.pos++ // skip '*'
	arraySizeBytes, err := p.readUntilCRLF()
	if err != nil {
		log.Fatal(err)
	}
	arraySize, err := strconv.Atoi(string(arraySizeBytes))
	if err != nil {
		log.Fatal(err)
	}

	command := &Command{
		Name: "",
		Args: make([]string, 0, arraySize),
	}

	for i := 0; i < arraySize; i++ {
		if p.pos >= p.len {
			log.Fatal("unexpected end of input while parsing array")
		}
		switch p.buf[p.pos] {
		case '*', '~', '>', '%': // Array
			subCommand := p.ParseArray()
			if i == 0 {
				command.Name = subCommand.Name
				command.Args = append(command.Args, subCommand.Args...)
			} else {
				// For nested arrays, add the name as an argument, then add all args
				command.Args = append(command.Args, subCommand.Name)
				command.Args = append(command.Args, subCommand.Args...)
			}
		case ':': // Integer
			integerStr := p.ParseInteger()
			if i == 0 {
				command.Name = integerStr.Name
			} else {
				command.Args = append(command.Args, integerStr.Name)
			}
		case '$': // Bulk string
			bulkStr := p.parseBulkString()
			if i == 0 {
				command.Name = bulkStr.Name
			} else {
				command.Args = append(command.Args, bulkStr.Name)
			}
		case '-', '+':
			simpleStr := p.parseSimpleString()
			if i == 0 {
				command.Name = simpleStr.Name
			} else {
				command.Args = append(command.Args, simpleStr.Name)
			}
		case '#': // Boolean
			boolStr := p.ParseBoolean()
			if i == 0 {
				command.Name = boolStr.Name
			} else {
				command.Args = append(command.Args, boolStr.Name)
			}
		case '!': // Blob error
			blobErr := p.parseBlobError()
			if i == 0 {
				command.Name = blobErr.Name
			} else {
				command.Args = append(command.Args, blobErr.Name)
			}
		case '_': // Null
			nullStr := p.ParseNull()
			if i == 0 {
				command.Name = nullStr.Name
			} else {
				command.Args = append(command.Args, nullStr.Name)
			}
		}
	}
	return command
}
func (p *StreamingParser) readUntilCRLF() ([]byte, error) {
	curr := p.pos
	for p.pos < p.len && p.buf[p.pos] != '\r' {
		p.pos++
	}
	if p.pos >= p.len {
		return nil, errors.New("unexpected end of input")
	}
	if p.pos+1 >= p.len || p.buf[p.pos+1] != '\n' {
		return nil, errors.New("no CRLF found")
	}
	result := p.buf[curr:p.pos]
	p.pos += 2 // skip \r\n
	return result, nil
}
func (p *StreamingParser) readSimpleString() ([]byte, error) {
	content, err := p.readUntilCRLF()
	if err != nil {
		return nil, err
	}
	return content, nil
}
func (p *StreamingParser) parseSimpleString() *Command {
	p.pos++ // skip '+' or '-'
	content, err := p.readSimpleString()
	if err != nil {
		log.Fatal(err)
	}
	// check if content is "+\r\n" or "-\r\n"
	if string(content) == "+\r\n" || string(content) == "-\r\n" {
		log.Fatal("invalid simple string value")
	}
	return &Command{
		Name: string(content),
		Args: []string{},
	}
}
func (p *StreamingParser) readBulkString() ([]byte, error) {
	p.pos++ // skip '$'
	lengthBytes, err := p.readUntilCRLF()
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return nil, err
	}

	if p.pos+length >= p.len {
		return nil, errors.New("bulk string content exceeds buffer")
	}

	content := make([]byte, length)
	copy(content, p.buf[p.pos:p.pos+length])
	p.pos += length

	// Skip CRLF if it exists
	if p.pos+1 < p.len && p.buf[p.pos] == '\r' && p.buf[p.pos+1] == '\n' {
		p.pos += 2
	}

	return content, nil
}
func (p *StreamingParser) parseBulkString() *Command {
	content, err := p.readBulkString()
	if err != nil {
		log.Fatal(err)
	}
	// check if content is "$0\r\n"
	if string(content) == "$0\r\n" {
		log.Fatal("invalid bulk string value")
	}
	return &Command{
		Name: string(content),
		Args: []string{},
	}
}
func (p *StreamingParser) readBlobError() ([]byte, error) {
	p.pos++ // skip '$'
	lengthBytes, err := p.readUntilCRLF()
	if err != nil {
		return nil, err
	}
	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return nil, err
	}

	if p.pos+length >= p.len {
		return nil, errors.New("blob error content exceeds buffer")
	}

	content := make([]byte, length)
	copy(content, p.buf[p.pos:p.pos+length])
	p.pos += length

	// Skip CRLF if it exists
	if p.pos+1 < p.len && p.buf[p.pos] == '\r' && p.buf[p.pos+1] == '\n' {
		p.pos += 2
	}

	return content, nil
}
func (p *StreamingParser) parseBlobError() *Command {
	content, err := p.readBulkString()
	if err != nil {
		log.Fatal(err)
	}
	// check if content is "$0\r\n"
	if string(content) == "!\r\n" {
		log.Fatal("invalid bulk string value")
	}
	return &Command{
		Name: string(content),
		Args: []string{},
	}
}
func (p *StreamingParser) skipCRLF() error {
	if p.pos+1 >= p.len {
		return io.EOF
	}

	if p.buf[p.pos] == '\r' && p.buf[p.pos+1] == '\n' {
		p.pos += 2
		return nil
	}

	return fmt.Errorf("expected CRLF at position %d", p.pos)
}
func (p *StreamingParser) ParseInteger() *Command {
	p.pos++ // skip ':'
	content, err := p.readUntilCRLF()
	if err != nil {
		log.Fatal(err)
	}
	//check if the content is valid integer
	if _, err := strconv.Atoi(string(content)); err != nil {
		log.Fatal("invalid integer value")
	}
	return &Command{
		Name: string(content),
		Args: []string{},
	}
}
// don't need the readUntilCRLF because we already read the #
func (p *StreamingParser) ParseBoolean() *Command {
	p.pos++ // skip '#'
	content, err := p.readUntilCRLF()
	if err != nil {
		log.Fatal(err)
	}
	// RESP3 boolean: #t\r\n or #f\r\n
	boolValue := string(content)
	// check if boolValue is "t" or "f"
	if boolValue != "t" && boolValue != "f" {
		log.Fatal("invalid boolean value")
	}
	return &Command{
		Name: boolValue,
		Args: []string{},
	}
}
func (p *StreamingParser) ParseNull() *Command {
	p.pos++ // skip '_'
	// Null is just _\r\n, no content
	err := p.skipCRLF()
	if err != nil {
		log.Fatal(err)
	}
	return &Command{
		Name: "null",
		Args: []string{},
	}
}
func NewStreamingParser(data []byte) *StreamingParser {
	return &StreamingParser{
		buf: data,
		pos: 0,
		len: len(data),
	}
}
