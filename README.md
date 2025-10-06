# YAKVS - Yet Another Key-Value Store

A Redis-compatible in-memory key-value store written in Go, implementing the RESP (Redis Serialization Protocol) for client communication.

## 🚀 Features

### ✅ Implemented Features

- **RESP Protocol Support**: Full implementation of Redis Serialization Protocol
  - Array commands (`*`)
  - Bulk strings (`$`)
  - Simple strings (`+`, `-`)
  - Integers (`:`)
  - Booleans (`#`)
  - Blob errors (`!`)
  - Null values (`_`)

- **Core Commands**:
  - `SET key value` - Set a key-value pair
  - `GET key` - Retrieve a value by key
  - `DEL key` - Delete a key
  - `EXISTS key` - Check if a key exists

- **Persistence**:
  - AOF (Append Only File) persistence
  - Automatic command logging for data-modifying operations
  - Recovery from AOF file on startup

- **Interactive Mode**:
  - Command-line interface with `>>` prompt
  - Automatic conversion of plain text commands to RESP format
  - Support for both RESP and plain text input

### 🏗️ Architecture

The project follows a modular architecture with clear separation of concerns:

```
YAKVS/
├── aof/                    # AOF persistence module
│   └── aof.go             # AOF file management
├── parser/                 # RESP protocol parser
│   ├── parser.go          # Streaming parser implementation
│   └── parser_test.go     # Comprehensive test suite
├── store/                  # Key-value storage
│   └── store.go           # In-memory store with interface
├── utils/                  # Utility functions
│   └── utils.go           # RESP conversion and validation
├── main.go                 # Main application entry point
├── execute.go              # Command execution engine
└── base.aof               # AOF persistence file
```

### 📦 Modules

#### AOF Module (`aof/`)
- **AOFManager**: Centralized AOF file operations
- **WriteCommand()**: Persist commands to AOF file
- **ReadAndExecuteCommands()**: Replay commands from AOF on startup
- **ShouldPersistCommand()**: Determine which commands to persist

#### Parser Module (`parser/`)
- **StreamingParser**: Efficient RESP protocol parsing
- **ParseCommand()**: Main parsing entry point
- **ParseArray()**: Handle RESP arrays
- **ParseBulkString()**: Handle RESP bulk strings
- **Comprehensive test coverage**: 100% test coverage for all parsing functions

#### Store Module (`store/`)
- **Store**: In-memory key-value storage
- **NewStore()**: Constructor for store instances
- **CRUD Operations**: GetValue, SetValue, DeleteValue, Exists
- **StoreInterface**: Interface for future extensibility

#### Utils Module (`utils/`)
- **ToRESP()**: Convert plain text commands to RESP format
- **IsRESPFormat()**: Detect RESP protocol format
- **PreprocessInput()**: Handle escape sequences (deprecated)

## 🚀 Getting Started

### Prerequisites

- Go 1.21.1 or later
- Git

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/shubhdevelop/YAKVS.git
   cd YAKVS
   ```

2. **Build the application**:
   ```bash
   go build -o YAKVS
   ```

3. **Run the application**:
   ```bash
   ./YAKVS
   ```

### Usage

#### Interactive Mode

Start the application and use the interactive prompt:

```bash
$ ./YAKVS
YAKVS
>> SET mykey "Hello World"
Parsing RESP command: *3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$11\r\nHello World\r\n
Executing command: &{Name:SET Args:[mykey Hello World]}
>> GET mykey
Parsing RESP command: *2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n
Executing command: &{Name:GET Args:[mykey]}
Hello World
>> DEL mykey
Parsing RESP command: *2\r\n$3\r\nDEL\r\n$5\r\nmykey\r\n
Executing command: &{Name:DEL Args:[mykey]}
>> EXISTS mykey
Parsing RESP command: *2\r\n$6\r\nEXISTS\r\n$5\r\nmykey\r\n
Executing command: &{Name:EXISTS Args:[mykey]}
false
>> exit
```

#### RESP Protocol Support

The application supports both plain text commands and native RESP protocol:

**Plain Text Input** (automatically converted):
```
>> SET key value
```

**Native RESP Input**:
```
>> *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
```

### Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run parser tests with verbose output
go test ./parser -v

# Run tests for specific module
go test ./store -v
```

## 🔧 Development

### Project Structure

- **Modular Design**: Each component is in its own package
- **Interface-Based**: Uses interfaces for extensibility
- **Test-Driven**: Comprehensive test coverage
- **Error Handling**: Proper error handling throughout

### Adding New Commands

1. **Add command logic in `execute.go`**:
   ```go
   case "NEWCOMMAND":
       // Implementation
   ```

2. **Update `utils/utils.go`** to support command conversion:
   ```go
   case "NEWCOMMAND":
       // Add to ToRESP() function
   ```

3. **Add to persistent commands** in `aof/aof.go`:
   ```go
   persistentCommands := map[string]bool{
       "NEWCOMMAND": true,
   }
   ```

### Code Quality

- **Linting**: No linting errors
- **Testing**: 100% test coverage for parser module
- **Documentation**: Comprehensive inline documentation
- **Error Handling**: Proper error propagation and handling

## 📊 Current Status

### ✅ Completed Features

- [x] RESP Protocol Parser
- [x] Core Key-Value Operations
- [x] AOF Persistence
- [x] Interactive Command Line Interface
- [x] Command Conversion (Plain Text → RESP)
- [x] Modular Architecture
- [x] Comprehensive Testing
- [x] Error Handling

### 🚧 In Progress

- [ ] Additional Redis Commands (HSET, HGET, LPUSH, etc.)
- [ ] Configuration Management
- [ ] Network Server Mode
- [ ] Clustering Support
- [ ] Memory Optimization

### 📋 Future Roadmap

- [ ] **Advanced Data Types**: Lists, Sets, Hashes, Sorted Sets
- [ ] **Expiration**: TTL support for keys
- [ ] **Persistence Options**: RDB snapshots, AOF rewriting
- [ ] **Network Protocol**: TCP server for remote connections
- [ ] **Replication**: Master-slave replication
- [ ] **Clustering**: Distributed key-value store
- [ ] **Performance**: Memory optimization, connection pooling
- [ ] **Monitoring**: Metrics and health checks

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Redis and its RESP protocol
- Built with Go's excellent concurrency primitives
- Test-driven development approach

---

**YAKVS** - A modern, modular key-value store implementation in Go 🚀