# TCP Chat Application (Go)

A multi-client TCP chat server where clients connect via TCP, register nickname, list users, send direct messages and broadcasts messages to groups.  

## How to run locally

### Prerequisites 
- Go 1.22+ (tested with Go 1.22.2)

### Start the server(terminal 1)
From the project root(where go.mod is)

```bash
go run ./cmd/server
```

### Start a client(terminal 2)
From the project root(where go.mod is)

```bash
go run ./cmd/client
```

Start many clients with different terminals to test 

## Commands

Commands are case-insensitive

### `/NCK <nickname>`
Set or change a nickname.

- Must start with alphabet optionally followed by alphanemeric or underscore. Up to 10 characters.

- Fail when the nickname is taken. 

Example:

/nck homer 

### /LST
Show list of registered nicknames. 

Example:

/LST

### `/MSG <recipients> <message>`
Send message to recipients. 

- Must set nickname first. 

Examples:

/MSG homer hello homer 

/msg homer,bart hello simpson!

### `/GRP <groupname> <users>`
Register a group for registered users.

- Group name must start with # followed by alphabet and optionally followed by alphanumeric or underscore character. Up to 11 characters.

Example:

/GRP #simpson homer,bart

/MSG #simpson hello simpson!

## Architecture
- Server listens on port 6666 and accepts clients.
- Each client connection starts a new goroutine (`HandleConnection`).
- Client goroutines send `Message` into a shared `messages` channel.
- Server retrieves `Message` from `messages` channel and replies through `Replies` channel which each client has.

This design avoids shared mutable state between goroutines.

Refer to Architecture.md for detailed internal structure.