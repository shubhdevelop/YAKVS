# YAKVS Command Reference

This document provides a comprehensive guide to all available commands in YAKVS (Yet Another Key-Value Store).

## Table of Contents

- [Getting Started](#getting-started)
- [Basic Commands](#basic-commands)
- [Numeric Commands](#numeric-commands)
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

## Numeric Commands

YAKVS supports numeric operations for integer values. These commands automatically handle numeric encoding and provide atomic operations.

### INCR

**Syntax:** `INCR key`

**Description:** Increments the integer value of a key by 1. If the key doesn't exist, it's initialized to 0 before incrementing.

**Arguments:**
- `key` (string): The key to increment

**Returns:**
- `:<new_value>` - The new value after incrementing

**Behavior:**
- **New Key**: Creates key with value 0, then increments to 1
- **Existing Key**: Increments the current value by 1
- **Atomic Operation**: Thread-safe increment operation
- **Integer Only**: Works only with integer values

**Example:**
```
>> INCR counter
:1
>> INCR counter
:2
>> INCR new_counter
:1
```

### DECR

**Syntax:** `DECR key`

**Description:** Decrements the integer value of a key by 1. If the key doesn't exist, it's initialized to 0 before decrementing.

**Arguments:**
- `key` (string): The key to decrement

**Returns:**
- `:<new_value>` - The new value after decrementing

**Behavior:**
- **New Key**: Creates key with value 0, then decrements to -1
- **Existing Key**: Decrements the current value by 1
- **Atomic Operation**: Thread-safe decrement operation
- **Integer Only**: Works only with integer values

**Example:**
```
>> DECR counter
:-1
>> DECR counter
:-2
>> DECR new_counter
:-1
```

### INCRBY

**Syntax:** `INCRBY key increment`

**Description:** Increments the integer value of a key by the specified amount. If the key doesn't exist, it's initialized to 0 before incrementing.

**Arguments:**
- `key` (string): The key to increment
- `increment` (integer): The amount to increment by (can be negative)

**Returns:**
- `:<new_value>` - The new value after incrementing

**Behavior:**
- **New Key**: Creates key with value 0, then adds the increment
- **Existing Key**: Adds the increment to the current value
- **Negative Increment**: Can be used to decrement (same as DECRBY with positive value)
- **Zero Increment**: Returns the same value (no change)
- **Atomic Operation**: Thread-safe increment operation

**Example:**
```
>> INCRBY counter 5
:5
>> INCRBY counter 3
:8
>> INCRBY counter -2
:6
>> INCRBY new_counter 10
:10
>> INCRBY counter 0
:6
```

### DECRBY

**Syntax:** `DECRBY key decrement`

**Description:** Decrements the integer value of a key by the specified amount. If the key doesn't exist, it's initialized to 0 before decrementing.

**Arguments:**
- `key` (string): The key to decrement
- `decrement` (integer): The amount to decrement by (can be negative)

**Returns:**
- `:<new_value>` - The new value after decrementing

**Behavior:**
- **New Key**: Creates key with value 0, then subtracts the decrement
- **Existing Key**: Subtracts the decrement from the current value
- **Negative Decrement**: Can be used to increment (same as INCRBY with positive value)
- **Zero Decrement**: Returns the same value (no change)
- **Atomic Operation**: Thread-safe decrement operation

**Example:**
```
>> DECRBY counter 3
:-3
>> DECRBY counter 2
:-5
>> DECRBY counter -4
:-1
>> DECRBY new_counter 7
:-7
>> DECRBY counter 0
:-1
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

### Working with Numeric Values

```bash
# Basic increment operations
>> INCR page_views
:1
>> INCR page_views
:2
>> INCRBY page_views 5
:7

# Basic decrement operations
>> DECR stock_count
:-1
>> DECR stock_count
:-2
>> DECRBY stock_count 3
:-5

# Working with existing numeric values
>> SET score 100
+OK
>> INCR score
:101
>> INCRBY score 50
:151
>> DECRBY score 25
:126

# Negative increments/decrements
>> INCRBY counter -10
:-10
>> DECRBY counter -5
:-5
>> INCRBY counter 15
:10

# Zero operations (no change)
>> INCRBY counter 0
:10
>> DECRBY counter 0
:10
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

### Numeric Operations Workflow

```bash
# Initialize counters
>> INCR user:1:visits
:1
>> INCRBY user:1:points 100
:100

# Track multiple metrics
>> INCRBY user:1:visits 4
:5
>> INCRBY user:1:points 50
:150
>> DECRBY user:1:points 25
:125

# Check final values
>> GET user:1:visits
5
>> GET user:1:points
125

# Reset a counter
>> DECRBY user:1:visits 5
:0
>> INCR user:1:visits
:1
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
   Error: INCR requires 1 argument (key)
   Error: DECR requires 1 argument (key)
   Error: INCRBY requires 2 arguments (key, value)
   Error: DECRBY requires 2 arguments (key, value)
   ```

2. **Invalid TTL/Timestamp:**
   ```
   Error parsing TTL: strconv.ParseInt: parsing "invalid": invalid syntax
   ```

3. **Invalid Numeric Values:**
   ```
   Error: INCRBY requires a valid integer value
   Error: DECRBY requires a valid integer value
   ```

4. **RESP Parsing Errors:**
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
- `INCR` commands are persisted
- `DECR` commands are persisted
- `INCRBY` commands are persisted
- `DECRBY` commands are persisted

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
