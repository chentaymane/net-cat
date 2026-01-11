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
	conn     net.Conn
	name     string
	messages []string
	mu       sync.Mutex
}

var (
	clients = make(map[net.Conn]*Client)
	mu      sync.Mutex
	history []string // store chat history
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

	// Send chat history to new client
	mu.Lock()
	for _, msg := range history {
		client.mu.Lock()
		fmt.Fprintln(client.conn, msg)
		client.mu.Unlock()
	}
	mu.Unlock()

	// Broadcast join message to all except self
	joinMsg := fmt.Sprintf("%s has joined our chat...", name)
	saveHistory(joinMsg)         // add to history
	broadcast(joinMsg, conn, false)

	for {
		msg, err := reader.ReadString('\n')
		
		if err != nil {
			leaveMsg := fmt.Sprintf("%s has left our chat...", name)
			fmt.Println(leaveMsg)
			saveHistory(leaveMsg)

			mu.Lock()
			delete(clients, conn)
			mu.Unlock()

			broadcast(leaveMsg, conn, true)
			return
		}
		msg = strings.TrimSpace(msg)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fullMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg)

		saveHistory(fullMsg)
		broadcast(fullMsg, conn, false)
	}
}

func saveHistory(msg string) {
	mu.Lock()
	history = append(history, msg)
	mu.Unlock()
}

// broadcast message to all clients
// isInfo=true → send to all (joins/leaves)
// isInfo=false → normal messages, skip sender
func broadcast(msg string, sender net.Conn, isInfo bool) {
	// Print to server log
	fmt.Println(msg) //

	mu.Lock()
	list := make([]*Client, 0, len(clients))
	for _, c := range clients {
		list = append(list, c)
	}
	mu.Unlock()
	for _, c := range list {
		if !isInfo && c.conn == sender {
			continue // skip sender for normal messages
		}
		c.mu.Lock()
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(c.conn, msg)
		c.mu.Unlock()
	}
}
