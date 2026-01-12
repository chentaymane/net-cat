package main

import (
	"bufio"
	"fmt"
	"log"
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
const MAX_CLIENT = 10 

const logo = `Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
`

func main() {
	listener, err := net.Listen("tcp", ":8080")
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
	if len(clients)>=MAX_CLIENT {
		conn.Write([]byte("max client now !"))
		return
	}
	reader := bufio.NewReader(conn)
	// Send logo
	conn.Write([]byte(logo))

	// Ask for name
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}

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

	// Broadcast join message
	joinMsg := fmt.Sprintf("%s has joined our chat...", name)
	log.Println(joinMsg)
	saveHistory(joinMsg)
	broadcastToOthers(joinMsg, conn)

	// Send initial prompt to this client
	sendPrompt(client)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Client disconnected
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			leaveMsg := fmt.Sprintf("[%s][%s]:%s has left our chat...", timestamp, name, name)
			saveHistory(leaveMsg)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			broadcastToAll(leaveMsg, nil)
			return
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			sendPrompt(client)
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fullMsg := fmt.Sprintf("[%s][%s]:%s", timestamp, name, msg)
		saveHistory(fullMsg)
		log.Println(fullMsg)
		// Broadcast to others
		broadcastToOthers(fullMsg, conn)
		
		// Send new prompt to sender
		sendPrompt(client)
	}
}

func saveHistory(msg string) {
	mu.Lock()
	history = append(history, msg)
	mu.Unlock()
}

func broadcastToOthers(msg string, senderConn net.Conn) {
	mu.Lock()
	list := make([]*Client, 0, len(clients))
	for _, c := range clients {
		list = append(list, c)
	}
	mu.Unlock()

	for _, c := range list {
		if c.conn == senderConn {
			continue
		}
		c.mu.Lock()
		// Print newline first to clear the current prompt line
		fmt.Fprintln(c.conn, "\n"+msg)
		sendPromptLocked(c)
		c.mu.Unlock()
	}
}

func broadcastToAll(msg string, skipConn net.Conn) {
	mu.Lock()
	list := make([]*Client, 0, len(clients))
	for _, c := range clients {
		list = append(list, c)
	}
	mu.Unlock()

	for _, c := range list {
		if c.conn == skipConn {
			continue
		}
		c.mu.Lock()
		// Print newline first to clear the current prompt line
		fmt.Fprintln(c.conn, "\n"+msg)
		sendPromptLocked(c)
		c.mu.Unlock()
	}
}

func sendPrompt(c *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sendPromptLocked(c)
}

func sendPromptLocked(c *Client) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(c.conn, "[%s][%s]:", timestamp, c.name)
}