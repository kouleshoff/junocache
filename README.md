# GO Redis-like in-memory cache

Learning GO by implementing a simplified caching server with

Desired features:
- Key-value storage with string, lists, dict support
- Per-key TTL
- Operations:
  - Get
  - Set
  - Update
  - Remove
  - Find all Keys
- Get i-th element on list
- Get value by key from dict
- Telnet-like API protocol

### Commands supported

The protocol is text-based over TCP/IP

Each command starts on a new line and consists of a verb, a cache key
and any optional arguments.
The SET command line should be followed by command body (any characters)
and terminated by two newlines to indicate the end of current command.

- `GET abc` (get string by key)
- `SET abc` (set string value for the specified key)
- `DEL abc` (delete value stored for the given key, if any)
- `exit` to stop the server

The response indicates whether the command succeeded
In case of a GET command, the value stored under specified key is written
to the response stream.

### Expiration of Keys

Optionally the command may have a Time-to-Live indication.
It applies to commands that modify data. This is possible
by appending "TTL n" to the first line of command, where n is in seconds.
After n seconds has elapsed from the insertion time, the value is no longer available.

For example: `SET abc TTL 10`

### Starting the Server

Use `go run ...` command to start the server on port 8080.
The port number may be changed directly in the source code, it is
specified by `const LISTEN_PORT ...`

```
go run server.go commands.go items.go cache.go
```

### Client connections

The service has been tested using telnet client.
Use empty line to close current client connection

```
-> telnet localhost 8080
Trying ::1...
Connected to localhost.
Escape character is '^]'.
GET a
-ERR GET with key a
SET a TTL 10
DHTSXKuDeGfIjsSB

+OK SET with key a
GET a
+OK GET with key a
DHTSXKuDeGfIjsSB

```
