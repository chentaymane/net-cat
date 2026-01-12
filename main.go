package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Client struct {
	conn net.Conn
	name string
	mu   sync.Mutex
}

var (
	clients = make(map[net.Conn]*Client)
	mu      sync.Mutex
	history []string
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Ask for name
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	client := &Client{conn: conn, name: name}

	// Add client
	mu.Lock()
	clients[conn] = client
	mu.Unlock()

	// Send chat history
	mu.Lock()
	for _, msg := range history {
		client.mu.Lock()
		fmt.Fprintln(client.conn, msg)
		client.mu.Unlock()
	}
	mu.Unlock()

	// Broadcast join message (without prompt)
	joinMsg := fmt.Sprintf("%s has joined the chat", name)
	saveHistory(joinMsg)
	broadcast(joinMsg, nil)
	sendPrompt(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Leave message (without prompt)
			leaveMsg := fmt.Sprintf("%s has left the chat", name)
			fmt.Println(leaveMsg)
			saveHistory(leaveMsg)

			mu.Lock()
			delete(clients, conn)
			mu.Unlock()

			broadcast(leaveMsg, nil)
			return
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			// ignore empty messages
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fullMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg)

		saveHistory(fullMsg)
		broadcast(fullMsg, conn)

		// Only send soft prompt to the client who just typed
		sendPrompt(client)
	}
}

func saveHistory(msg string) {
	mu.Lock()
	history = append(history, msg)
	mu.Unlock()
}

func broadcast(msg string, skipSender net.Conn) {
	fmt.Println(msg) // server log

	mu.Lock()
	list := make([]*Client, 0, len(clients))
	for _, c := range clients {
		list = append(list, c)
	}
	mu.Unlock()

	for _, c := range list {
		if c.conn == skipSender {
			continue
		}
		c.mu.Lock()
		fmt.Fprintln(c.conn, msg)
		c.mu.Unlock()
	}
}

// send a soft prompt ONLY to this client
func sendPrompt(c *Client) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	c.mu.Lock()
	fmt.Fprint(c.conn, "[", timestamp, "][", c.name, "]: ")
	c.mu.Unlock()
}
