package main

import (
  "bufio"
  "fmt"
  "log"
  "net"
  "os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6666")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	stdinScanner := bufio.NewScanner(os.Stdin)
	connScanner  := bufio.NewScanner(conn)

	// Goroutine for receiving message
	go func() {
		for connScanner.Scan() {
			line := connScanner.Text()
			fmt.Println(line)
		}
		log.Println("connection closed by server")
		os.Exit(0) 
	}()

	// Send message to goroutine
	for {
		if !stdinScanner.Scan() {
			break
		}
		text := stdinScanner.Text()
		fmt.Fprintf(conn, "%s\n", text)
	}

	log.Println("client exiting")
}