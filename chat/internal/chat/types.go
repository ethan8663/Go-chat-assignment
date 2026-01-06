package chat

// Logging the status of client's command.
type Result string

const (
	Success    Result = "Success log"
	Fail       Result = "Fail log"
	NewMessage Result = "New message"
)

// Server to send result of the client's command back to client.
type Reply struct {
	Status Result
	Detail string
}

// Sent to server via Message struct.
type Client struct {
	Nickname string     // Server can change.
	Replies  chan Reply // Server can reply back.
}

// Client(goroutine) send message struct to server.
type Message struct {
	MType     Command
	Client    *Client 
	Detail    string
	Recipient []string // slice of nicknames.
}

// Available client's commands
type Command string

const (
    CmdNick    Command = "/NCK"
	CmdList    Command = "/LST"
    CmdMsg     Command = "/MSG"
	CmdGrp     Command = "/GRP"
	ClientExit Command = "/EXIT" // When client disconnects, automatically sent to server.
)

