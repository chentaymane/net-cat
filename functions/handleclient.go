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
	mu.Lock()
	if len(clients) >= MAX_CLIENT {
		conn.Write([]byte("\x1b[38;5;227mmax client now !\x1b[0m"))
		return
	}
	mu.Unlock()
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
		if validName(name) {
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
	defer DeleteClient(client)
	// Add client
	mu.Lock()
	clients[name] = client

	// Send chat history

	for _, msg := range history {
		fmt.Fprintln(client.conn, msg)
	}
	mu.Unlock()

	// Broadcast join message
	joinMsg := fmt.Sprintf("\x1b[38;5;46m%s has joined our chat...\x1b[0m", name)
	log.Println(joinMsg)
	saveHistory(joinMsg)
	broadcast(joinMsg, conn)

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
			delete(clients, name)
			mu.Unlock()
			broadcast(leaveMsg, nil)
			return
		}

		msg = strings.TrimSpace(msg)
		if !validMsg(msg) {
			sendPrompt(client)
			continue
		}
		tag := false

		for _, r := range clients {
			if r != client && strings.HasPrefix(msg, "@"+r.name) && strings.TrimSpace(msg[len("@"+r.name):]) != "" {
				Tag(client, r, msg)
				tag = true
				break
			}

		}

		if tag {
			continue
		}

		if TYPE, check := isAPrompt(msg); check {
			switch TYPE {
			case "rename":
				newName := strings.TrimPrefix(msg, "/rename ")
				if validName(newName) {
					RenameClient(client, newName)
				} else {
					client.conn.Write([]byte("\x1b[38;5;199minvalid name\n\x1b[0m"))
				}
				NewNameMsg := fmt.Sprintf(
					"\x1b[1;38;5;196m[%s] \x1b[38;5;173mhas changed his name to \x1b[38;5;196m[%s]\x1b[0m",
					name, newName,
				)
				log.Println(NewNameMsg)
				saveHistory(NewNameMsg)
				broadcast(NewNameMsg, conn)
			case "users":
				printUsers(conn)
			}
			sendPrompt(client)
			continue
		}

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fullMsg := fmt.Sprintf("\x1b[1;37m[%s][%s]:%s\x1b[0m", timestamp, name, "\x1b[38;5;51m"+msg+"\x1b[0m")
		saveHistory(fullMsg)
		log.Printf("\x1b[1;37m[%s]:%s\x1b[0m\n", name, msg)
		// Broadcast to others
		broadcast(fullMsg, conn)

		// Send new prompt to sender
		sendPrompt(client)
	}
}
