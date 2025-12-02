# Redis Protocol (RESP) - Complete Documentation

## Table of Contents
1. [Overview](#overview)
2. [Protocol Versions](#protocol-versions)
3. [Network Layer](#network-layer)
4. [Request-Response Model](#request-response-model)
5. [RESP Data Types](#resp-data-types)
6. [Command Format](#command-format)
7. [Pipelining](#pipelining)
8. [Client Handshake](#client-handshake)
9. [Inline Commands](#inline-commands)

---

## Overview

**RESP (REdis Serialization Protocol)** is the wire protocol used by Redis clients to communicate with the Redis server.

### Design Goals
- **Simple to implement**: Straightforward parsing logic
- **Fast to parse**: Uses prefixed lengths, comparable performance to binary protocols
- **Human readable**: Easy to debug with tools like telnet
- **Binary-safe**: Uses prefixed length for bulk data transfer

### Core Principles
- First byte determines the data type
- All protocol parts terminated with `\r\n` (CRLF)
- Prefixed lengths avoid scanning for special characters
- No escaping or quoting needed for bulk data

---

## Protocol Versions

### RESP2 (Standard since Redis 2.0)
- The most widely supported version
- Default protocol for Redis connections
- 5 core data types: Simple Strings, Errors, Integers, Bulk Strings, Arrays

### RESP3 (Available since Redis 6.0)
- Superset of RESP2 with additional types
- Adds: Booleans, Doubles, Maps, Sets, Pushes, and more
- Better semantic clarity (dedicated Null type)
- Opt-in via `HELLO` command
- Push-based notifications support

**Default**: Connections start in RESP2 mode. Clients can upgrade to RESP3 using the `HELLO` command.

---

## Network Layer

### Connection Details
- **Transport**: TCP connections (or Unix sockets)
- **Default Port**: 6379
- **Protocol**: Stream-oriented, request-response model

```
Client ──TCP Socket (port 6379)──> Redis Server
```

---

## Request-Response Model

### Standard Flow
1. Client sends command as RESP Array of Bulk Strings
2. Server processes the command
3. Server sends reply as appropriate RESP type

### Exceptions to Request-Response
1. **Pipelining**: Client sends multiple commands without waiting for individual replies
2. **Pub/Sub**: Protocol becomes push-based after subscription
3. **MONITOR**: Connection switches to ad-hoc push mode
4. **Protected Mode**: Server sends `-DENIED` and terminates connection
5. **RESP3 Pushes**: Server can send out-of-band data at any time

---

## RESP Data Types

### Type Summary Table

| Type | First Byte | Version | Category | Description |
|------|-----------|---------|----------|-------------|
| Simple Strings | `+` | RESP2 | Simple | Short non-binary strings |
| Simple Errors | `-` | RESP2 | Simple | Error messages |
| Integers | `:` | RESP2 | Simple | 64-bit signed integers |
| Bulk Strings | `$` | RESP2 | Bulk | Binary-safe strings |
| Arrays | `*` | RESP2 | Aggregate | Ordered collections |
| Nulls | `_` | RESP3 | Simple | Null value |
| Booleans | `#` | RESP3 | Simple | True/False |
| Doubles | `,` | RESP3 | Simple | Floating point |
| Big Numbers | `(` | RESP3 | Simple | Arbitrary precision integers |
| Bulk Errors | `!` | RESP3 | Bulk | Binary-safe errors |
| Verbatim Strings | `=` | RESP3 | Bulk | Strings with encoding hint |
| Maps | `%` | RESP3 | Aggregate | Key-value dictionaries |
| Sets | `~` | RESP3 | Aggregate | Unordered unique elements |
| Pushes | `>` | RESP3 | Aggregate | Out-of-band data |
| Attributes | `|` | RESP3 | Aggregate | Metadata annotations |

---

## RESP2 Data Types (Core)

### 1. Simple Strings

**Format**: `+<string>\r\n`

**Constraints**:
- Cannot contain CR (`\r`) or LF (`\n`) characters
- Terminated by CRLF
- Used for short, non-binary responses

**Examples**:
```
+OK\r\n
+PONG\r\n
```

---

### 2. Simple Errors

**Format**: `-<error message>\r\n`

**Constraints**:
- Similar to simple strings but indicates error condition
- First uppercase word is the error prefix
- Clients should treat as exceptions

**Common Error Prefixes**:
- `ERR` - Generic error
- `WRONGTYPE` - Operation against wrong data type
- `NOAUTH` - Authentication required
- `NOPROTO` - Protocol version not supported

**Examples**:
```
-ERR unknown command 'asdf'\r\n
-WRONGTYPE Operation against a key holding the wrong kind of value\r\n
```

---

### 3. Integers

**Format**: `:[<+|->]<value>\r\n`

**Constraints**:
- Signed 64-bit integer
- Base-10 representation
- Optional sign (+ or -)
- CRLF terminated

**Examples**:
```
:0\r\n
:1000\r\n
:-42\r\n
:+15\r\n
```

**Common Uses**:
- `INCR`, `DECR` - Counter values
- `LLEN`, `SCARD` - Collection sizes
- `EXISTS`, `SISMEMBER` - Boolean results (1=true, 0=false)
- `LASTSAVE` - UNIX timestamps

---

### 4. Bulk Strings

**Format**: `$<length>\r\n<data>\r\n`

**Constraints**:
- Binary-safe (can contain any byte sequence)
- Length in bytes (not characters)
- Default maximum: 512 MB (configurable via `proto-max-bulk-len`)
- Length-prefixed to avoid special character scanning

**Examples**:

Empty string:
```
$0\r\n\r\n
```

String "hello":
```
$5\r\nhello\r\n
```

Binary data with null bytes:
```
$10\r\nhe\x00llo\x00wo\r\n
```

**Null Bulk String** (RESP2 representation of NULL):
```
$-1\r\n
```
Represents a non-existent value (e.g., `GET` on missing key).

**Important**: 
- Length is in bytes, not UTF-8 characters
- Multi-byte characters count as multiple bytes
- `$3\r\n€\r\n` is incorrect if € is encoded as 3 bytes in UTF-8

---

### 5. Arrays

**Format**: `*<number-of-elements>\r\n<element-1>...<element-n>`

**Constraints**:
- Can contain any RESP type (including nested arrays)
- Elements can be of different types
- Supports unlimited nesting depth

**Examples**:

Empty array:
```
*0\r\n
```

Array of two bulk strings ["hello", "world"]:
```
*2\r\n
$5\r\nhello\r\n
$5\r\nworld\r\n
```

Array of three integers [1, 2, 3]:
```
*3\r\n
:1\r\n
:2\r\n
:3\r\n
```

Mixed-type array [1, 2, 3, 4, "hello"]:
```
*5\r\n
:1\r\n
:2\r\n
:3\r\n
:4\r\n
$5\r\nhello\r\n
```

Nested arrays [[1, 2, 3], ["Hello", Error("World")]]:
```
*2\r\n
*3\r\n
:1\r\n
:2\r\n
:3\r\n
*2\r\n
+Hello\r\n
-World\r\n
```

**Null Array** (RESP2):
```
*-1\r\n
```
Represents a null value (e.g., `BLPOP` timeout).

**Arrays with Null Elements**:
```
*3\r\n
$5\r\nhello\r\n
$-1\r\n
$5\r\nworld\r\n
```
Represents: `["hello", nil, "world"]`

Used by commands like `SORT` with `GET` pattern when keys are missing.

---

## RESP3 Additional Types

### 6. Nulls

**Format**: `_\r\n`

**Purpose**: Single dedicated null type, replacing RESP2's null bulk strings (`$-1\r\n`) and null arrays (`*-1\r\n`).

**Example**:
```
_\r\n
```

---

### 7. Booleans

**Format**: `#<t|f>\r\n`

**Examples**:
```
#t\r\n  (true)
#f\r\n  (false)
```

**Purpose**: Explicit boolean type instead of using integers (1/0).

---

### 8. Doubles

**Format**: `,[<+|->]<integral>[.<fractional>][<E|e>[sign]<exponent>]\r\n`

**Examples**:
```
,1.23\r\n
,10\r\n
,1.23e-4\r\n
,inf\r\n
,-inf\r\n
,nan\r\n
```

**Purpose**: Floating-point numbers with support for infinity and NaN.

---

### 9. Big Numbers

**Format**: `([<+|->]<number>\r\n`

**Purpose**: Integer values outside the 64-bit signed range.

**Examples**:
```
(3492890328409238509324850943850943825024385\r\n
(-12345678901234567890\r\n
```

**Note**: Cannot include fractional parts.

---

### 10. Bulk Errors

**Format**: `!<length>\r\n<error>\r\n`

**Purpose**: Error messages that can contain binary data or exceed simple string limits.

**Example**:
```
!21\r\n
SYNTAX invalid syntax\r\n
```

**Convention**: Error begins with uppercase word (error prefix).

---

### 11. Verbatim Strings

**Format**: `=<length>\r\n<encoding>:<data>\r\n`

**Constraints**:
- First 3 bytes after length indicate encoding
- Colon (`:`) separates encoding from data
- Total length includes encoding prefix

**Example**:
```
=15\r\n
txt:Some string\r\n
```

**Purpose**: Provides hint about data encoding (e.g., `txt` for plain text, `mkd` for markdown).

**Use Case**: Commands like `INFO` return formatted output that should be displayed as-is.

---

### 12. Maps

**Format**: `%<number-of-entries>\r\n<key-1><value-1>...<key-n><value-n>`

**Purpose**: Dictionary/hash type with key-value pairs.

**Example**:
```
%2\r\n
+first\r\n
:1\r\n
+second\r\n
:2\r\n
```
Represents: `{"first": 1, "second": 2}`

**RESP2 Alternative**: Flat array with alternating keys and values:
```
*4\r\n
$5\r\nfirst\r\n
:1\r\n
$6\r\nsecond\r\n
:2\r\n
```

**Note**: Both keys and values can be any RESP type.

---

### 13. Sets

**Format**: `~<number-of-elements>\r\n<element-1>...<element-n>`

**Purpose**: Unordered collection of unique elements.

**Example**:
```
~3\r\n
+apple\r\n
+banana\r\n
+orange\r\n
```

**Note**: Client libraries should return native set type if available.

---

### 14. Pushes

**Format**: `><number-of-elements>\r\n<element-1>...<element-n>`

**Purpose**: Out-of-band data sent by server (not in response to a command).

**Characteristics**:
- Can appear at any time
- Used for Pub/Sub messages
- Server-initiated notifications
- Similar structure to arrays but different semantic meaning

**Example** (Pub/Sub message):
```
>3\r\n
$7\r\nmessage\r\n
$7\r\nchannel\r\n
$5\r\nhello\r\n
```

**Important**: Pushes never appear inside other types (e.g., not in middle of an array).

---

### 15. Attributes

**Format**: `|<number-of-entries>\r\n<key-1><value-1>...<key-n><value-n>`

**Purpose**: Metadata that augments the following reply (not part of the reply itself).

**Structure**: Identical to Map type but uses `|` prefix.

**Example**:
```
|1\r\n
+key-popularity\r\n
%2\r\n
$1\r\na\r\n
,0.1923\r\n
$1\r\nb\r\n
,0.0012\r\n
*2\r\n
:2039123\r\n
:9543892\r\n
```

The actual reply is `[2039123, 9543892]`.
The attribute provides metadata: `{key-popularity: {a: 0.1923, b: 0.0012}}`.

**Position**: Attributes appear immediately before the data they describe.

**Use Case**: Key popularity, TTL info, command statistics, etc.

---

## Command Format

### Sending Commands to Redis

**All commands are sent as RESP Arrays of Bulk Strings.**

**General Format**:
```
*<argument-count>\r\n
$<len-arg1>\r\n<arg1>\r\n
$<len-arg2>\r\n<arg2>\r\n
...
```

### Examples

**Command**: `SET mykey myvalue`

**Encoded**:
```
*3\r\n
$3\r\n
SET\r\n
$5\r\n
mykey\r\n
$7\r\n
myvalue\r\n
```

Breakdown:
- `*3` - Array of 3 elements
- `$3\r\nSET` - Bulk string "SET" (3 bytes)
- `$5\r\nmykey` - Bulk string "mykey" (5 bytes)
- `$7\r\nmyvalue` - Bulk string "myvalue" (7 bytes)

---

**Command**: `GET mykey`

**Encoded**:
```
*2\r\n
$3\r\n
GET\r\n
$5\r\n
mykey\r\n
```

---

**Command**: `LLEN mylist`

**Encoded**:
```
*2\r\n
$4\r\n
LLEN\r\n
$6\r\n
mylist\r\n
```

---

**Command**: `INCR counter`

**Encoded**:
```
*2\r\n
$4\r\n
INCR\r\n
$7\r\n
counter\r\n
```

---

### Server Response Examples

**Command**: `SET mykey myvalue`

**Response**:
```
+OK\r\n
```

---

**Command**: `GET mykey`

**Response** (key exists):
```
$7\r\n
myvalue\r\n
```

**Response** (key doesn't exist):
```
$-1\r\n
```

---

**Command**: `LLEN mylist`

**Response**:
```
:48293\r\n
```

---

**Command**: `LRANGE mylist 0 2`

**Response**:
```
*3\r\n
$5\r\n
first\r\n
$6\r\n
second\r\n
$5\r\n
third\r\n
```

---

## Pipelining

### Concept
Pipelining allows clients to send multiple commands without waiting for individual replies, then read all replies in sequence.

### Benefits
- Reduced network round-trips (RTT)
- Improved throughput
- Lower overall latency for multiple commands

### How It Works
1. Client sends multiple commands in succession
2. Server processes each command in order
3. Client reads all replies in the same order

### Example

**Client Sends** (3 commands pipelined):
```
*3\r\n$3\r\nSET\r\n$4\r\nkey1\r\n$6\r\nvalue1\r\n
*3\r\n$3\r\nSET\r\n$4\r\nkey2\r\n$6\r\nvalue2\r\n
*2\r\n$3\r\nGET\r\n$4\r\nkey1\r\n
```

**Server Responds** (3 replies):
```
+OK\r\n
+OK\r\n
$6\r\nvalue1\r\n
```

### Important Notes
- Replies are returned in the same order as commands
- Server processes commands sequentially
- Client must read the exact number of replies sent
- Error in one command doesn't stop processing of subsequent commands

---

## Client Handshake

### HELLO Command

New connections should begin with the `HELLO` command to:
1. Negotiate protocol version
2. Optionally authenticate
3. Receive server information

**Format**:
```
HELLO <protocol-version> [AUTH <username> <password>] [SETNAME <clientname>]
```

### Examples

**Upgrade to RESP3**:
```
Client: HELLO 3
Server: %7\r\n
        +server\r\n
        +redis\r\n
        +version\r\n
        +7.0.0\r\n
        +proto\r\n
        :3\r\n
        ...
```

**RESP2-only server**:
```
Client: HELLO 3
Server: -ERR unknown command 'HELLO'\r\n
(Connection remains in RESP2)
```

**Unsupported version**:
```
Client: HELLO 4
Server: -NOPROTO sorry, this protocol version is not supported\r\n
```

**With authentication**:
```
Client: HELLO 3 AUTH default mypassword
Server: (Map reply with server info)
```

### Response Fields (RESP3)

**Mandatory fields**:
- `server` - Server software name (e.g., "redis")
- `version` - Server version string
- `proto` - Highest supported RESP version

**Redis-specific fields**:
- `id` - Connection ID
- `mode` - "standalone", "sentinel", or "cluster"
- `role` - "master" or "replica"
- `modules` - Array of loaded module names

### Default Behavior
- Connections start in RESP2 mode by default
- `HELLO` command upgrades to requested version
- Future Redis versions may change the default

---

## Inline Commands

### Purpose
Allows sending commands in a human-friendly format using tools like `telnet` when `redis-cli` is unavailable.

### Format
Space-separated arguments (no RESP encoding).

### How Redis Detects Inline Commands
If the first byte is NOT `*` (array indicator), Redis parses the command as inline format.

### Examples

**Using telnet**:
```
$ telnet localhost 6379
> PING
+PONG
> EXISTS somekey
:0
> SET mykey myvalue
+OK
> GET mykey
$7
myvalue
```

### Limitations
- Not suitable for binary data
- Less efficient than RESP format
- Primarily for debugging and manual testing
- Production clients should use proper RESP encoding

---

## Special Considerations

### Binary Safety
- RESP is fully binary-safe for bulk strings
- Use bulk strings (not simple strings) for any binary data
- Length prefix eliminates need for escaping

### Character Encoding
- Protocol itself is encoding-agnostic
- Length in bulk strings is bytes, not characters
- Client/server must agree on encoding (typically UTF-8)

### Performance
- Prefixed lengths enable single-pass parsing
- No need to scan for special characters
- Comparable to binary protocol performance
- Minimal overhead for human readability

### Compatibility
- RESP2 is universal and always supported
- RESP3 is opt-in and backward compatible
- Clients can detect server capabilities via `HELLO`
- Old clients work with new servers (in RESP2 mode)

### Size Limits
- Bulk strings: 512 MB default (configurable)
- Arrays: No practical limit on element count
- Integer: 64-bit signed (-2^63 to 2^63-1)
- Big numbers: Arbitrary precision in RESP3

---

## RESP3 Streaming Types (Advanced)

### Overview
RESP3 introduces streaming capabilities for situations where the size of data is not known in advance. This is particularly useful for:
- Real-time search results
- Modular Redis extensions returning dynamic data
- Replication streams
- Large datasets with unknown length

**Note**: Streaming strings and streaming aggregates were excluded from Redis 6.0's initial RESP3 support but are part of the full RESP3 specification.

### Streamed Strings

**Purpose**: Send strings whose length is initially unknown.

**Format**: `$?<\r\n><chunk-length>\r\n<data>\r\n...<terminator>`

**Characteristics**:
- Use `$?` instead of `$<length>` to indicate streaming mode
- Each chunk has format: `;<length>\r\n<data>\r\n`
- Terminated with `;0\r\n` (zero-length chunk)

**Example**:
```
$?\r\n
;5\r\nhello\r\n
;6\r\n world\r\n
;0\r\n
```

This represents the string "hello world" sent in chunks.

### Streamed Aggregates

**Purpose**: Send arrays, sets, or maps without knowing the final count.

**Format**: Use `?` instead of count, terminate with `.` (period)

**Streamed Array**:
```
*?\r\n
+element1\r\n
+element2\r\n
:123\r\n
.\r\n
```

**Streamed Set**:
```
~?\r\n
+apple\r\n
+banana\r\n
.\r\n
```

**Streamed Map**:
```
%?\r\n
+key1\r\n
:100\r\n
+key2\r\n
:200\r\n
.\r\n
```

**Terminator**: Single period (`.`) followed by CRLF marks the end.

**Nesting**: Streamed aggregates can contain other RESP types, including other streamed types.

---

## Protocol Behaviors and Edge Cases

### 1. Connection State and Mode Changes

**Standard Mode**: Request-response
**Special Modes**:
- **Pub/Sub**: After `SUBSCRIBE`, connection becomes push-based
- **MONITOR**: Connection streams all server commands
- **Blocking Commands**: Commands like `BLPOP` may timeout and return null

**Important**: In RESP2, these mode changes are implicit (based on connection state). In RESP3, Push types make this explicit.

### 2. Error Handling Patterns

**Protocol Errors vs Redis Errors**:
- **Protocol Errors**: Malformed RESP data (connection should close)
- **Redis Errors**: Command errors (e.g., wrong type, syntax error)

**Error Recovery**:
- Simple/Bulk errors are part of normal operation
- Client should raise exceptions for Redis errors
- Protocol parse errors should terminate connection

### 3. Partial Read Handling

**Buffering Strategy**:
- RESP requires reading complete tokens
- Length-prefixed types allow pre-allocation
- Must handle incomplete reads (especially over slow networks)

**Example Scenario**:
```
Receive: "*2\r\n$5\r\nhe"
Status: Incomplete (need 3 more bytes + CRLF)
Action: Buffer and wait for more data
```

### 4. Maximum Sizes and Limits

**Default Redis Limits**:
- Bulk strings: 512 MB (`proto-max-bulk-len`)
- Array elements: No hard limit (memory constrained)
- Integer range: -2^63 to 2^63-1
- Nesting depth: No specification limit (implementation dependent)

**Configuration**:
```
# redis.conf
proto-max-bulk-len 536870912  # 512MB in bytes
```

### 5. Unicode and Encoding

**Protocol Level**:
- RESP itself is encoding-agnostic
- Bulk strings are pure binary (byte sequences)
- Simple strings should be ASCII-safe

**Application Level**:
- Client and server must agree on encoding (usually UTF-8)
- Length is always in bytes, not characters
- Multi-byte characters require careful handling

**Example**:
- String "€" in UTF-8 is 3 bytes
- Correct: `$3\r\n€\r\n` (3 bytes)
- Incorrect: `$1\r\n€\r\n` (claims 1 byte)

### 6. Empty vs Null vs Missing

**Empty String**: `$0\r\n\r\n` (string of zero length)
**Null (RESP2)**: `$-1\r\n` (non-existent value)
**Null (RESP3)**: `_\r\n` (explicit null)
**Empty Array**: `*0\r\n` (array with zero elements)
**Null Array (RESP2)**: `*-1\r\n` (non-existent array)

**Semantic Differences**:
- Empty: Value exists but has no content
- Null: Value does not exist
- Commands like `GET` return null for missing keys
- Commands like `LRANGE` return empty array for empty lists

### 7. Transaction Responses

**MULTI/EXEC Pattern**:
```
Client: MULTI
Server: +OK\r\n

Client: SET key1 value1
Server: +QUEUED\r\n

Client: SET key2 value2
Server: +QUEUED\r\n

Client: EXEC
Server: *2\r\n
        +OK\r\n
        +OK\r\n
```

**Structure**: `EXEC` returns array containing individual command replies.

### 8. Pub/Sub Protocol Behavior

**RESP2 Pub/Sub** (implicit push mode):
```
*3\r\n
$7\r\nmessage\r\n
$7\r\nchannel\r\n
$5\r\nhello\r\n
```

**RESP3 Pub/Sub** (explicit push type):
```
>3\r\n
$7\r\nmessage\r\n
$7\r\nchannel\r\n
$5\r\nhello\r\n
```

**Key Difference**: RESP3 uses `>` prefix making it distinguishable from command replies.

### 9. Protected Mode

When Redis is in protected mode and receives a connection from a non-loopback address:
```
Server: -DENIED Redis is running in protected mode...\r\n
(Connection terminated immediately)
```

### 10. Security: SSRF Prevention

Redis implements protections against Server-Side Request Forgery attacks, particularly when malformed HTTP requests are sent to Redis servers. If Redis detects suspicious patterns (like HTTP headers), it may disconnect.

**Example** (sending HTTP to Redis):
```
GET /key HTTP/1.1
Host: localhost
```
Redis may detect "Host:" and treat this as a potential attack, disconnecting the client.

---

## Advanced Topics

### Protocol Negotiation

**Backward Compatibility Flow**:
1. Client connects (starts in RESP2 mode)
2. Client sends `HELLO 3` to request RESP3
3. Server responds:
   - Success: Map with server info (now in RESP3 mode)
   - Failure: `-ERR` or `-NOPROTO` (remains in RESP2)
4. Client adapts based on response

### Cluster Protocol

**Important**: The RESP protocol documented here is for client-server communication only.

**Redis Cluster**: Uses a completely different binary protocol for node-to-node communication (not RESP).

**Sentinel Mode**: Uses RESP but with specific command patterns for high availability.

### Client-Side Caching (RESP3 Feature)

RESP3 enables server-push invalidation messages for client-side caching:
```
>2\r\n
$10\r\ninvalidate\r\n
*1\r\n
$4\r\nkey1\r\n
```

Server pushes invalidation when cached keys are modified.

### Blocking Operations

Commands like `BLPOP`, `BRPOP`, `BRPOPLPUSH`:
- Block connection until data available or timeout
- On timeout, return null array: `*-1\r\n` (RESP2) or `_\r\n` (RESP3)
- Client must handle long-running operations

### Command Pipelining Best Practices

**Benefits**:
- Reduces RTT from N × RTT to 1 × RTT
- Increases throughput significantly

**Limitations**:
- No transactions (commands execute independently)
- Must read all replies in order
- Error in one command doesn't affect others

**When Not to Pipeline**:
- Commands depend on previous results
- Need immediate error feedback
- Memory constraints (buffering many responses)

---

## Implementation Checklist

### Minimal RESP2 Database Implementation

**Must Support**:
- ✓ Parse: Simple Strings, Errors, Integers, Bulk Strings, Arrays
- ✓ Encode: Arrays of Bulk Strings (commands)
- ✓ Respond: Appropriate RESP2 types per command
- ✓ Handle: Null bulk strings (`$-1`)
- ✓ Handle: Null arrays (`*-1`)
- ✓ Handle: Empty strings and arrays
- ✓ Support: Binary-safe bulk strings
- ✓ Support: TCP socket connections
- ✓ Support: Pipelining (queue multiple commands)

**Recommended**:
- ✓ Implement: Inline command parsing (for telnet debugging)
- ✓ Implement: Error prefixes (ERR, WRONGTYPE, etc.)
- ✓ Handle: Connection timeouts
- ✓ Handle: Partial reads/writes
- ✓ Validate: Length limits (prevent DoS)

### Optional RESP3 Features

**Core RESP3**:
- ✓ Parse: All RESP3 types (Null, Boolean, Double, etc.)
- ✓ Implement: `HELLO` command handshake
- ✓ Support: Protocol version negotiation
- ✓ Support: Attributes (metadata)
- ✓ Support: Push types (out-of-band data)

**Advanced RESP3**:
- ✓ Implement: Streamed strings
- ✓ Implement: Streamed aggregates
- ✓ Support: Client-side caching
- ✓ Support: Push-based Pub/Sub

---

## Reference: Complete Type Reference

| Type | Format | Example | Notes |
|------|--------|---------|-------|
| Simple String | `+str\r\n` | `+OK\r\n` | No CR/LF in string |
| Error | `-err\r\n` | `-ERR unknown\r\n` | First word is error prefix |
| Integer | `:num\r\n` | `:1000\r\n` | 64-bit signed |
| Bulk String | `$len\r\ndata\r\n` | `$5\r\nhello\r\n` | Binary-safe |
| Array | `*count\r\n...` | `*2\r\n:1\r\n:2\r\n` | Can nest |
| Null | `_\r\n` | `_\r\n` | RESP3 only |
| Boolean | `#t\r\n` or `#f\r\n` | `#t\r\n` | RESP3 only |
| Double | `,num\r\n` | `,1.23\r\n` | RESP3 only |
| Big Number | `(num\r\n` | `(123...\r\n` | RESP3 only |
| Bulk Error | `!len\r\nerr\r\n` | `!5\r\nERROR\r\n` | RESP3 only |
| Verbatim | `=len\r\nenc:data\r\n` | `=11\r\ntxt:hello\r\n` | RESP3 only |
| Map | `%count\r\n...` | `%1\r\n+a\r\n:1\r\n` | RESP3 only |
| Set | `~count\r\n...` | `~2\r\n+a\r\n+b\r\n` | RESP3 only |
| Push | `>count\r\n...` | `>3\r\n...` | RESP3 only |
| Attribute | `|count\r\n...` | `|1\r\n+key\r\n:1\r\n` | RESP3 only |
| Streamed String | `$?\r\n;len\r\ndata\r\n;0\r\n` | See above | RESP3, optional |
| Streamed Array | `*?\r\n...\r\n.\r\n` | See above | RESP3, optional |
| Streamed Set | `~?\r\n...\r\n.\r\n` | See above | RESP3, optional |
| Streamed Map | `%?\r\n...\r\n.\r\n` | See above | RESP3, optional |

---

## Quick Reference: Common Patterns

### Reading a Command
```
*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
```
1. Read `*` → Array type
2. Read `3` → 3 elements
3. For each element, read type and data

### Sending Success
```
+OK\r\n
```

### Sending an Error
```
-ERR unknown command\r\n
```

### Returning a String Value
```
$11\r\nhello world\r\n
```

### Returning Null (key not found)
```
$-1\r\n              (RESP2)
_\r\n                (RESP3)
```

### Returning Multiple Values
```
*3\r\n
$5\r\nvalue1\r\n
$5\r\nvalue2\r\n
$5\r\nvalue3\r\n
```

### Returning Integer
```
:42\r\n
```