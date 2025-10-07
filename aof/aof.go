package aof

import (
	"fmt"
	"io"
	"os"

	"github.com/shubhdevelop/YAKVS/parser"
)

type AOFManager struct {
	writeFile *os.File
	readFile  *os.File
	filename  string
}

func NewAOFManager(filename string) *AOFManager {
	return &AOFManager{
		filename: filename,
	}
}


func (aof *AOFManager) Initialize() error {
	// Open file for writing (AOF - Append Only File)
	writeFile, err := os.OpenFile(aof.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening write file: %v", err)
	}
	aof.writeFile = writeFile

	// Open file for reading
	readFile, err := os.Open(aof.filename)
	if err != nil {
		// If file doesn't exist yet, that's okay - we'll create it when we write
		if !os.IsNotExist(err) {
			return fmt.Errorf("error opening read file: %v", err)
		}
		aof.readFile = nil
		return nil
	}
	aof.readFile = readFile
	return nil
}


func (aof *AOFManager) Close() error {
	var err error
	if aof.writeFile != nil {
		if closeErr := aof.writeFile.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if aof.readFile != nil {
		if closeErr := aof.readFile.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

func (aof *AOFManager) WriteCommand(command string) error {
	if aof.writeFile == nil {
		return fmt.Errorf("write file not initialized")
	}
	
	_, err := aof.writeFile.WriteString(command)
	if err != nil {
		return fmt.Errorf("failed to write to AOF file: %v", err)
	}
	
	// Flush to ensure data is written to disk
	return aof.writeFile.Sync()
}

func (aof *AOFManager) ReadAndExecuteCommands(executeFunc func(*parser.Command)) error {
	if aof.readFile == nil {
		fmt.Println("No AOF file found, starting fresh.")
		return nil
	}

	fmt.Println("Reading from AOF file:")
	
	// Read the entire file content
	fileContent, err := io.ReadAll(aof.readFile)
	if err != nil {
		return fmt.Errorf("error reading AOF file: %v", err)
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
		executeFunc(command)
	}
	
	return nil
}

func (aof *AOFManager) ShouldPersistCommand(commandName string) bool {
	// Only persist commands that modify data
	persistentCommands := map[string]bool{
		"SET": true,
		"DEL": true,
		"EXPIRE": true,
		"EXPIREAT": true,
		// Add more commands that modify data as needed
	}
	return persistentCommands[commandName]
}

func (aof *AOFManager) GetWriteFile() *os.File {
	return aof.writeFile
}

func (aof *AOFManager) GetReadFile() *os.File {
	return aof.readFile
}
