package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"
	// "log"
)

func HandleConnection(conn net.Conn, messages chan Message) {
  defer conn.Close()

  scanner := bufio.NewScanner(conn)

  replies := make(chan Reply)

  // ex. #simpsons : ["homer", "bart"]
  groups := make(map[string][]string)

  // default nickname is ""
  client := &Client{Nickname: "", Replies: replies}

  go func() {
	for reply := range client.Replies {
		fmt.Fprintf(conn, "%s\n", handleReply(reply))
	}
  }()

  for scanner.Scan() {
	userCommand := strings.TrimSpace(scanner.Text())
	if userCommand == "" {
		continue
	}

	parts := strings.SplitN(userCommand, " ", 2)
	first := strings.ToUpper(parts[0])
	cmd   := Command(first)

	var rest string
	if len(parts) == 2 {
		rest = parts[1]
	}
	
	switch cmd {
	case CmdNick:
		if err := validateNck(rest); err != nil {
			fmt.Fprintf(conn, "%s\n", err.Error())
			continue
		}

		nicknameParts := strings.Fields(rest)
		nickname := nicknameParts[0]

		messages <- Message{
			MType : CmdNick,
			Client: client,
			Detail: nickname,
		}

	case CmdList: 
		messages <- Message {
			MType : CmdList,
			Client: client,
		}

	case CmdMsg:
		if err := validateMsg(rest, client); err != nil {
			fmt.Fprintf(conn, "%s\n", err.Error())
			continue
		}

		parts            := strings.SplitN(rest, " ", 2)
		recipientsString := parts[0]
		detail           := parts[1]

		recipientsArray := strings.Split(recipientsString, ",")

		var recipients []string

		for _, rec := range recipientsArray {
			rec = strings.TrimSpace(rec)
			if rec == "" {
				continue
			}

			if strings.HasPrefix(rec, "#") {
				groupName := rec 
				members, ok := groups[groupName]
				if !ok {
					fmt.Fprintf(conn, "[Fail log] No such group: %s\n", groupName)
					continue
				}
				recipients = append(recipients, members...)
			} else {
				recipients = append(recipients, rec)
			}
		}

		messages <- Message {
			MType    : CmdMsg,
			Client   : client,
			Detail   : detail,
			Recipient: recipients,
		}

	// Group command is not sent to server.  
	case CmdGrp:
		if err := validateGrp(rest, client); err != nil {
			fmt.Fprintf(conn, "[Fail log] %s\n", err.Error())
			continue
		}

		parts       := strings.SplitN(rest, " ", 2)
		groupName   := parts[0]
		groupMember := parts[1]

		groupMemberArray := strings.Split(groupMember, ",")

		groups[groupName] = groupMemberArray

		client.Replies <- Reply {
			Status: Success,
			Detail: "Group is registered",
		}

	default:
		fmt.Fprintf(conn, "[Fail log] Invalid command\n")
	}

	}

	// After client disconnects, remove client nickname.
	messages <- Message {
		MType : ClientExit,
		Detail: client.Nickname,
	}

	fmt.Println("connection closed:", conn)
}

// ---- private functions for validating command ----

func validateNck(rest string) error {
	atLeastOneArg := validNumOfArg(1)
	max10Length   := validNumOfLength(10)

	if !atLeastOneArg(rest) {
		return fmt.Errorf("Number of argument should be at least 1.")
	}

	parts    := strings.Fields(rest)
	nickname := parts[0]
	if !max10Length(nickname) {
		return fmt.Errorf("Can not be longer than 10")
	}

	if !startWithAlpha(nickname) {
		return fmt.Errorf("Should start with alphabet")
	}

	if !containsAlphaNumOrUnderScore(nickname) {
		return fmt.Errorf("Should only contain alphabet, number, or underscore")
	}

	return nil
}

func validateMsg(rest string, client *Client) error {
	atLeastTwoArg := validNumOfArg(2)

	if !atLeastTwoArg(rest) {
		return fmt.Errorf("Number of argument should be at least 2")
	}

	if client.Nickname == "" {
		return fmt.Errorf("Should set nickname first")
	}

	return nil
}

func validateGrp(rest string, client *Client) error {
	atLeastTwoArg := validNumOfArg(2)
	max11Length   := validNumOfLength(11)

	if !atLeastTwoArg(rest) {
		return fmt.Errorf("Number of argument should be at least 2")
	}

	if client.Nickname == "" {
		return fmt.Errorf("Should set nickname first")
	}

	parts     := strings.Fields(rest)
	groupName := parts[0]
	if !strings.HasPrefix(groupName, "#") {
		return fmt.Errorf("Group name should start with #")
	}

	nameWithoutHash := groupName[1:]
	if !max11Length(nameWithoutHash) {
		return fmt.Errorf("Can not be longer than 11")
	}

	if !containsAlphaNumOrUnderScore(nameWithoutHash) {
		return fmt.Errorf("Should only contain alphabet, number, or underscore")
	}

	return nil
}

// ---- helper function ----

func validNumOfArg(num int) func(string) bool {
	return func(rest string) bool {
		if strings.TrimSpace(rest) == "" {
			return false
		}
		parts := strings.Fields(rest)
		return len(parts) >= num
	}
}

func validNumOfLength(num int) func(string) bool {
	return func(name string) bool {
		return len(name) <= num
	}
}

func startWithAlpha(name string) bool {
	if name == "" {
		return false
	}
	r := []rune(name)[0] 
    return unicode.IsLetter(r)
}

func containsAlphaNumOrUnderScore(name string) bool {
	for _, r := range name {
        if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
            return false
        }
    }
    return true
}

func handleReply(reply Reply) string {
	return "[" + string(reply.Status) + "] " + reply.Detail
}