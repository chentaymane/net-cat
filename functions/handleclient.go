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

	// Check max clients (protected)

	reader := bufio.NewReader(conn)

	// Send logo to client
	conn.Write([]byte(logo))

	// Ask for client name
	conn.Write([]byte("\x1b[1;37m[ENTER YOUR NAME]: \x1b[0m"))
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)

	// Validate name
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

	// Create client
	client := &Client{conn: conn, name: name}
	TYPE := Add(client)
	if TYPE != "" {
		if TYPE == "name" {
			conn.Write([]byte("\x1b[38;5;199mName already exists.\n\x1b[0m"))
		} else {
			conn.Write([]byte("\x1b[38;5;199mServer full, try again later.\n\x1b[0m"))
			return
		}
	}
	defer Remove(client) // Ensure cleanup on exit

	// Send chat history to new client
	loadHistory(client)

	// Broadcast join message
	joinMsg := fmt.Sprintf("\x1b[38;5;46m%s has joined our chat...\x1b[0m", name)
	log.Println(joinMsg)
	saveHistory(joinMsg)
	broadcast(joinMsg, conn)

	// Send initial prompt
	sendPrompt(client)

	for {
		// Read message from client
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

		// Ignore invalid messages
		if !validMsg(msg) {
			sendPrompt(client)
			continue
		}

		// Handle tagged/private messages (@user)
		tag := false
		for _, r := range clients {
			if r != client &&
				strings.HasPrefix(msg, "@"+r.name) &&
				strings.TrimSpace(msg[len("@"+r.name):]) != "" {

				Tag(client, r, msg)
				tag = true
				break
			}
		}

		if tag {
			continue
		}

		// Handle commands (/rename, /users)
		if TYPE, ok := isAPrompt(msg); ok {
			switch TYPE {
			case "rename":
				newName := strings.TrimPrefix(msg, "/rename ")
				if validName(newName) {
					Rename(client, newName)
				} else {
					client.conn.Write([]byte("\x1b[38;5;199minvalid name\n\x1b[0m"))
				}

				newNameMsg := fmt.Sprintf(
					"\x1b[1;38;5;196m[%s] \x1b[38;5;173mhas changed his name to \x1b[38;5;196m[%s]\x1b[0m",
					name, newName,
				)
				log.Println(newNameMsg)
				saveHistory(newNameMsg)
				broadcast(newNameMsg, conn)

			case "users":
				printUsers(conn)
			}
		} else {

			// Normal chat message
			timestamp := time.Now().Format("2006-01-02 15:04:05")
			fullMsg := fmt.Sprintf(
				"\x1b[1;37m[%s][%s]:%s\x1b[0m",
				timestamp,
				name,
				"\x1b[38;5;51m"+msg+"\x1b[0m",
			)

			saveHistory(fullMsg)
			log.Printf("\x1b[1;37m[%s]:%s\x1b[0m\n", name, msg)

			// Broadcast message to others
			broadcast(fullMsg, conn)
		}
		// Send prompt back to sender
		sendPrompt(client)
	}
}
