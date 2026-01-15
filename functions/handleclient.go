package functions

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()
	if len(clients) >= MAX_CLIENT {
		conn.Write([]byte("\x1b[38;5;227mmax client now !\x1b[0m"))
		return
	}
	reader := bufio.NewReader(conn)
	// Send logo
	conn.Write([]byte(logo))

	// Ask for name
	conn.Write([]byte("\x1b[1;37m[ENTER YOUR NAME]: \x1b[0m"))
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	for {
		if validName(name) && name != "" {
			break
		}
		conn.Write([]byte("\x1b[38;5;199mName invalid or already exists !\n[ENTER NEW NAME]: \x1b[0m"))
		name, err = reader.ReadString('\n')
		if err != nil {
			return
		}
		name = strings.TrimSpace(name)
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
	joinMsg := fmt.Sprintf("\x1b[38;5;46m %s has joined our chat...\x1b[0m", name)
	log.Println(joinMsg)
	saveHistory(joinMsg)
	broadcastToOthers(joinMsg, conn)

	// Send initial prompt to this client
	sendPrompt(client)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Client disconnected
			leaveMsg := fmt.Sprintf("\x1b[38;5;197m%s has left our chat...\x1b[0m", name)
			log.Println(leaveMsg)
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
		fullMsg := fmt.Sprintf("\x1b[1;37m[%s][%s]:%s\x1b[0m", timestamp, name, "\x1b[38;5;51m"+msg+"\x1b[0m")
		saveHistory(fullMsg)
		log.Println(fmt.Sprintf("\x1b[1;37m[%s]:%s\x1b[0m", name, msg))
		// Broadcast to others
		broadcastToOthers(fullMsg, conn)

		// Send new prompt to sender
		sendPrompt(client)
	}
}
