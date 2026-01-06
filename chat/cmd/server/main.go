package main

import (
	"log"
	"net"
	"chat/internal/chat"
	"strings"
)

// Server starts goroutine when new client connects.
// Listens to messages channel.
func main() {
	ln, err := net.Listen("tcp", ":6666")
	if err != nil {
		log.Fatal(err)
	}

	messages  := make(chan chat.Message)
	nicknames := make(map[string]*chat.Client)

	// goroutine for accepting new client.
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("accept error:", err)
				continue
			}

			go chat.HandleConnection(conn, messages)
		}
	}()

	// Keeps receiving message from client goroutine.
	for message := range messages {
		switch message.MType {
		case chat.CmdNick:
			nickname  := message.Detail
			_, exists := nicknames[nickname]
			if exists {
				message.Client.Replies <- chat.Reply{Status: chat.Fail, Detail: nickname + " already exists",}
			} else {
				oldNickname := message.Client.Nickname

				// If client changes nickname, nicknames map has to delete previous nickname.
				if oldNickname != "" {
					delete(nicknames, oldNickname)
				}
				
				message.Client.Nickname = nickname
				nicknames[nickname]     = message.Client

				message.Client.Replies <- chat.Reply{Status: chat.Success, Detail: "changed nickname to " + nickname,}
			}

		case chat.CmdList:
			names := make([]string, 0, len(nicknames))

			for n := range nicknames {
				names = append(names, n)
			}

			if len(names) == 0 {
				message.Client.Replies <- chat.Reply {
					Status: chat.Success, Detail: "There is no occupied nicknames",
				}
			} else {
				line := "Occupied nicknames: " + strings.Join(names, ", ")
				message.Client.Replies <- chat.Reply {
					Status: chat.Success, Detail: line,
				}
			}

		case chat.CmdMsg:
			for _, rec := range message.Recipient {
				target, ok := nicknames[rec]
				if !ok {
					message.Client.Replies <- chat.Reply {
						Status: chat.Fail, 
						Detail: "No user nickname as " + rec,
					}
					continue
				}
				target.Replies <- chat.Reply {
					Status: chat.NewMessage, 
					Detail: message.Client.Nickname + ": " + message.Detail,
				}
			}

		case chat.ClientExit:
			nickname := message.Detail
			delete(nicknames, nickname)
			log.Printf("%s is deleted from server", nickname)
		}
	}
}