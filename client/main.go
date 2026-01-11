package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(os.Stdin)

	// Server prompt for name
	prompt, _ := serverReader.ReadString(':')
	fmt.Print(prompt)

	// Enter name
	scanner.Scan()
	name := scanner.Text()
	conn.Write([]byte(name + "\n"))
	
	// Goroutine: receive messages live
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	go func() {
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("\nDisconnected from server")
				os.Exit(0)
			}
			msg = strings.TrimSpace(msg)
			// print messages from others or join/leave info
			fmt.Println("\n" + msg)
			fmt.Print("[", timestamp, "][", name, "]: ") // reprint prompt
		}
	}()

	// Main goroutine: read input and send to server
	for {
		fmt.Print("[", timestamp, "][", name, "]: ")
		if !scanner.Scan() {
			return
		}
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Exiting...")
			return
		}
		conn.Write([]byte(text + "\n"))
	}
}
