package functions

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// Save a message to chat history (thread-safe)
func saveHistory(msg string) {
	mu.Lock()
	defer mu.Unlock()
	history = append(history, msg)
}

// Public function: safe prompt sending
func sendPrompt(c *Client) {
	mu.Lock()
	defer mu.Unlock()
	sendPromptLocked(c)
}

// Internal function: assumes mutex is already locked
func sendPromptLocked(c *Client) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(
		c.conn,
		"\x1b[1;37m[%s][%s]:\x1b[0m",
		timestamp,
		c.name,
	)
}

// Validate client name
func validName(name string) bool {
	name = strings.TrimSpace(name)

	// Empty or too long
	if name == "" || len(name) > 10 {
		return false
	}

	// Only lowercase letters and digits
	for _, r := range strings.ToLower(name) {
		if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') {
			return false
		}
	}

	// Check uniqueness
	mu.Lock()
	_, exists := clients[name]
	mu.Unlock()

	return !exists
}

// Validate chat message
func validMsg(msg string) bool {
	msg = strings.TrimSpace(msg)

	// Empty or too long
	if msg == "" || len(msg) > 100 {
		return false
	}

	// Allowed characters only
	for _, r := range msg {
		if !strings.ContainsRune(allowedChars, r) {
			return false
		}
	}

	return true
}

// Remove client from map
func DeleteClient(c *Client) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, c.name)
}

// Rename a client safely
func RenameClient(c *Client, newName string) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, c.name)
	c.name = newName
	clients[newName] = c
}

// Print all connected users
func printUsers(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range clients {
		fmt.Fprintln(conn, "\x1b[1;38;5;226m"+c.name+"\x1b[0m")
	}
}

// Handle private/tagged message (@user)
func Tag(sender *Client, receiver *Client, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Extract private message body
	privateMsg := strings.TrimSpace(msg[len("@"+receiver.name):])

	formatted := fmt.Sprintf(
		"\x1b[1;37m[%s][%s]:%s\x1b[0m",
		timestamp,
		sender.name,
		"\x1b[44m"+privateMsg+"\x1b[0m",
	)

	// Send to receiver
	fmt.Fprintln(receiver.conn, "\n"+formatted)
	sendPrompt(receiver)

	// Log server-side
	log.Printf("\x1b[1;37m[%s]:%s\x1b[0m", sender.name, msg)

	// Restore sender prompt
	sendPrompt(sender)
}

// Detect commands
func isAPrompt(str string) (string, bool) {
	str = strings.TrimSpace(str)

	if strings.HasPrefix(str, "/rename ") {
		return "rename", true
	}

	if str == "/users" {
		return "users", true
	}

	return "", false
}
