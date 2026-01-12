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
Server 
data:
messages   chat.Message channel:    receive message from client goroutine.
nicknames  map[string]*chat.Client: store nickname and corresponding client as key-value.

Server is listening to port to accept new client.
Server is listening to messages channel to get message from client goroutine. 

HandleConnection goroutine(server starts this when client connects)
data:
replies chan Reply         : channel that is stored inside of client struct. 
groups  map[string][]string: store group name (start with #) 
                             and corresponding slice of nicknames.
client  *Client            : has nickname and channel for receiving reply from server.

Goroutine is listening to client's command.
Goroutine is listening to replies channel to receive the result from server. 

## Types 
Result: string for logging the status of client's command.
It has 3 constants; 
Success    : command was successful.
Fail       : command was unsuccessful.
NewMessage : another user sent message to this client. 

Reply struct: used by server to send result of the client's command back to client.
Status: status of result.
Detail: content.

Client struct: used by client. This struct is sent inside of Message struct when the client send command to goroutine and the goroutine forward it to server.
Nickname: nickname of client. Server can modify it when necessary.
Replies : server constructs Reply and put into this channel.

Message struct: used by client. Client sends command to command to goroutine. Goroutine consturcts Message struct based on client's command. Put Message into messages channel so that server can pick up. 
MType    : type of message.
Client   : pointer of client.
Detail   : content of message.
Recipient: slice of nicknames.

Command: used by client. Sort client's command into one of constants.
CmdNick   : command for setting nickname.
CmdList   : command for viewing registered nicknames.
CmdMsg    : command for sending message.
CmdGrp    : command for registering a group.
ClientExit: this is not a command that client can send. It is used when client disconnects. 



