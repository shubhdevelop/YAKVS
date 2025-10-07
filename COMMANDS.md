# YAKVS Command Reference

This document provides a comprehensive guide to all available commands in YAKVS (Yet Another Key-Value Store).

## Table of Contents

- [Getting Started](#getting-started)
- [Basic Commands](#basic-commands)
- [TTL and Expiration Commands](#ttl-and-expiration-commands)
- [Command Syntax](#command-syntax)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [RESP Protocol Support](#resp-protocol-support)

## Getting Started

YAKVS supports both interactive mode and RESP protocol. Start the application with:

```bash
./YAKVS
```

You'll see the prompt:
```
YAKVS
>> 
```

## Basic Commands

### SET

**Syntax:** `SET key value`

**Description:** Sets a key-value pair in the store.

**Arguments:**
- `key` (string): The key to set
- `value` (string): The value to associate with the key

**Returns:** No output on success

**Example:**
```
>> SET mykey "Hello World"
>> SET username "john_doe"
>> SET counter "100"
```

### GET

**Syntax:** `GET key`

**Description:** Retrieves the value associated with a key.

**Arguments:**
- `key` (string): The key to retrieve

**Returns:** 
- The value if the key exists
- `nil` if the key doesn't exist

**Example:**
```
>> GET mykey
Hello World
>> GET nonexistent
nil
```

### DEL

**Syntax:** `DEL key`

**Description:** Deletes a key and its associated value from the store.

**Arguments:**
- `key` (string): The key to delete

**Returns:** No output on success

**Example:**
```
>> DEL mykey
>> DEL username
```

### EXISTS

**Syntax:** `EXISTS key`

**Description:** Checks if a key exists in the store.

**Arguments:**
- `key` (string): The key to check

**Returns:**
- `true` if the key exists
- `false` if the key doesn't exist

**Example:**
```
>> EXISTS mykey
true
>> EXISTS nonexistent
false
```

## TTL and Expiration Commands

### TTL

**Syntax:** `TTL key`

**Description:** Returns the remaining time-to-live (TTL) for a key. Automatically cleans up expired keys.

**Arguments:**
- `key` (string): The key to check TTL for

**Returns:**
- `:<remaining_seconds>` - Remaining seconds until expiration
- `:-1` - Key exists but has no expiration set
- `:-2` - Key doesn't exist or has expired (automatically cleaned up)

**Behavior:**
- **Automatic Cleanup**: Expired keys are automatically deleted when accessed
- **Dynamic Calculation**: Returns actual remaining seconds, not Unix timestamp
- **Memory Efficient**: Expired keys are removed from storage immediately

**Example:**
```
>> TTL mykey
:3599
>> TTL no_expiry_key
:-1
>> TTL expired_key
:-2
```

### EXPIRE

**Syntax:** `EXPIRE key seconds`

**Description:** Sets an expiration time (TTL) for a key in seconds from now.

**Arguments:**
- `key` (string): The key to set expiration for
- `seconds` (integer): Number of seconds until expiration

**Returns:**
- `+OK` - Successfully set expiration
- `:0` - Key doesn't exist

**Example:**
```
>> EXPIRE mykey 3600
+OK
>> EXPIRE session_token 1800
+OK
>> EXPIRE nonexistent 3600
:0
```

### EXPIREAT

**Syntax:** `EXPIREAT key timestamp`

**Description:** Sets an expiration time for a key using a Unix timestamp.

**Arguments:**
- `key` (string): The key to set expiration for
- `timestamp` (integer): Unix timestamp when the key should expire

**Returns:**
- `+OK` - Successfully set expiration
- `:0` - Key doesn't exist

**Example:**
```
>> EXPIREAT mykey 1735689600
+OK
>> EXPIREAT session_token 1735693200
+OK
>> EXPIREAT nonexistent 1735689600
:0
```

### PERSIST

**Syntax:** `PERSIST key`

**Description:** Removes the expiration from a key, making it persistent (no TTL).

**Arguments:**
- `key` (string): The key to make persistent

**Returns:**
- `:1` - Successfully removed expiration
- `:0` - Key doesn't exist or has no expiration set

**Behavior:**
- **Removes TTL**: Key becomes persistent and won't expire
- **Safe Operation**: Works on keys with or without existing TTL
- **Memory Efficient**: Removes key from expiry tracking

**Example:**
```
>> SET mykey "value"
+OK
>> EXPIRE mykey 3600
+OK
>> TTL mykey
:3599
>> PERSIST mykey
:1
>> TTL mykey
:-1
>> PERSIST nonexistent
:0
```

## Command Syntax

### Interactive Mode

YAKVS supports plain text commands in interactive mode:

```
>> COMMAND arg1 arg2 arg3
```

### RESP Protocol

YAKVS also supports native RESP (Redis Serialization Protocol) format:

```
>> *3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$11\r\nHello World\r\n
```

## Examples

### Basic Key-Value Operations

```bash
# Set a key-value pair
>> SET name "John Doe"

# Retrieve the value
>> GET name
John Doe

# Check if key exists
>> EXISTS name
true

# Delete the key
>> DEL name

# Check if key still exists
>> EXISTS name
false
```

### Working with TTL

```bash
# Set a key with value
>> SET session "abc123"
+OK

# Set expiration to 1 hour (3600 seconds)
>> EXPIRE session 3600
+OK

# Check TTL (returns remaining seconds)
>> TTL session
:3599

# Set expiration to specific timestamp
>> EXPIREAT session 1735689600
+OK

# Check TTL again (shows remaining time)
>> TTL session
:7200

# Check TTL for key with no expiration
>> TTL session
:-1

# Check TTL for expired key (automatically cleaned up)
>> TTL expired_key
:-2

# Make a key persistent (remove TTL)
>> PERSIST session
:1
>> TTL session
:-1
```

### Complete Workflow Example

```bash
# Start with a clean store
>> SET user:1 "Alice"
>> SET user:2 "Bob"

# Check existence
>> EXISTS user:1
true
>> EXISTS user:3
false

# Set expiration for user:1
>> EXPIRE user:1 7200

# Check TTL
>> TTL user:1
7199

# Get values
>> GET user:1
Alice
>> GET user:2
Bob

# Make user:1 persistent (remove TTL)
>> PERSIST user:1
:1
>> TTL user:1
:-1

# Delete user:2
>> DEL user:2

# Verify deletion
>> EXISTS user:2
false
>> GET user:2
nil
```

## Automatic Expiration Behavior

YAKVS implements intelligent automatic expiration with the following features:

### 🔄 Automatic Cleanup
- **Lazy Expiration**: Expired keys are automatically deleted when accessed
- **Memory Efficient**: No background processes needed - cleanup happens on-demand
- **Immediate Removal**: Expired keys are removed from storage instantly

### ⏱️ TTL Response Values
- `:<positive_number>` - Remaining seconds until expiration
- `:-1` - Key exists but has no expiration set
- `:-2` - Key doesn't exist or has expired (automatically cleaned up)

### 🎯 Best Practices
- **Check TTL before operations**: Use `TTL` to verify key status
- **Handle expired keys gracefully**: Expect `-2` responses for expired keys
- **Set reasonable expiration times**: Avoid very short TTLs for frequently accessed keys

### Example: Automatic Expiration
```bash
# Set a key with short expiration
>> SET temp "data"
+OK
>> EXPIRE temp 1
+OK

# Wait for expiration, then check
>> TTL temp
:-2

# Key is automatically deleted
>> GET temp
$-1
>> EXISTS temp
:0
```

## Error Handling

### Common Error Messages

1. **Insufficient Arguments:**
   ```
   Error: SET requires 2 arguments (key, value)
   Error: GET requires 1 argument (key)
   Error: DEL requires 1 argument (key)
   Error: EXISTS requires 1 argument (key)
   Error: TTL requires 1 argument (key)
   Error: EXPIRE requires 2 arguments (key, ttl)
   Error: EXPIREAT requires 2 arguments (key, timestamp)
   ```

2. **Invalid TTL/Timestamp:**
   ```
   Error parsing TTL: strconv.ParseInt: parsing "invalid": invalid syntax
   ```

3. **RESP Parsing Errors:**
   ```
   Error converting to RESP: [error details]
   Error parsing RESP command: [error details]
   ```

### Error Handling Best Practices

- Always provide the correct number of arguments
- Use valid integer values for TTL and timestamps
- Check command syntax before execution
- Handle nil values appropriately in your applications

## RESP Protocol Support

YAKVS implements the RESP (Redis Serialization Protocol) for compatibility with Redis clients.

### Supported RESP Types

- **Arrays** (`*`): For command arrays
- **Bulk Strings** (`$`): For command arguments
- **Simple Strings** (`+`, `-`): For simple responses
- **Integers** (`:`): For numeric values
- **Booleans** (`#`): For true/false values
- **Blob Errors** (`!`): For error messages
- **Null** (`_`): For null values

### RESP Examples

**Plain Text Input:**
```
>> SET key value
```

**Equivalent RESP Input:**
```
>> *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
```

**GET Command in RESP:**
```
>> *2\r\n$3\r\nGET\r\n$3\r\nkey\r\n
```

## Special Commands

### Interactive Commands

- `exit` - Exit the application
- `clear` - Clear the screen
- `help` - Show help (if implemented)

### Persistence

YAKVS automatically persists data-modifying commands to the AOF (Append Only File) for durability:

- `SET` commands are persisted
- `DEL` commands are persisted
- `EXPIRE` commands are persisted
- `EXPIREAT` commands are persisted

Read-only commands (`GET`, `EXISTS`, `TTL`) are not persisted.

## Performance Notes

- All operations are performed in-memory for maximum speed
- TTL expiration is checked on access
- AOF persistence may impact write performance
- Memory usage grows with the number of stored keys

## Troubleshooting

### Common Issues

1. **Command not recognized:** Ensure you're using the correct command syntax
2. **TTL not working:** Check that the key exists and TTL is set correctly
3. **Persistence issues:** Verify AOF file permissions and disk space
4. **Memory issues:** Monitor memory usage with large datasets

### Debug Mode

Enable debug output by checking the console for:
- Command parsing information
- RESP conversion details
- Execution status messages

---

For more information about YAKVS architecture and development, see the main [README.md](README.md) file.
