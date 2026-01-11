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

	// Broadcast join message to all except self
	broadcast(fmt.Sprintf("%s has joined our chat...", name), conn, false)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			broadcast(fmt.Sprintf("%s has left our chat...", name), conn, true)
			return
		}
		msg = strings.TrimSpace(msg)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		broadcast(fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg), conn, false)
	}
}

// broadcast message to all clients
// isInfo=true → send to all (joins/leaves)
// isInfo=false → normal messages, skip sender
func broadcast(msg string, sender net.Conn, isInfo bool) {
	mu.Lock()
	list := make([]*Client, 0, len(clients))
	for _, c := range clients {
		list = append(list, c)
	}
	mu.Unlock()

	for _, c := range list {
		if !isInfo && c.conn == sender {
			continue // skip sender
		}
		c.mu.Lock()
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(c.conn, msg)
		c.mu.Unlock()
	}
}
