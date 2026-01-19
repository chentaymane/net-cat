package functions

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func saveHistory(msg string) {
	mu.Lock()
	defer mu.Unlock()
	history = append(history, msg)
}

func sendPrompt(c *Client) {
	mu.Lock()
	defer mu.Unlock()
	sendPromptLocked(c)
}

func sendPromptLocked(c *Client) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(c.conn, "\x1b[1;37m[%s][%s]:\x1b[0m", timestamp, c.name)
}

func validName(name string) bool {
	if strings.TrimSpace(name) == "" || len(name) > 10 {
		return false
	}
	for _, r := range strings.ToLower(name) {
		if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') {
			return false
		}
	}
	mu.Lock()
	_, exists := clients[name]
	mu.Unlock()
	return !exists
}

func validMsg(msg string) bool {
	mu.Lock()
	defer mu.Unlock()

	if strings.TrimSpace(msg) == "" || len(msg) > 100 {
		return false
	}

	for _, r := range msg {
		if !strings.ContainsRune(allowedChars, r) {
			return false
		}
	}

	return true
}

func DeleteClient(c *Client) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, c.name)
}

func RenameClient(c *Client, newName string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, c.name)
	c.name = newName
	clients[newName] = c
}

func printUsers(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	for _, c := range clients {
		fmt.Fprintln(conn, "\x1b[1;38;5;226m"+c.name+"\x1b[0m")
	}
}

func Tag(sender *Client, receiver *Client, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	PRVmsg := strings.TrimSpace(msg[len("@"+receiver.name):])

	formatted := fmt.Sprintf("\x1b[1;37m[%s][%s]:%s\x1b[0m", timestamp, sender.name, "\x1b[44m"+PRVmsg+"\x1b[0m")
	fmt.Fprintln(receiver.conn, "\n"+formatted)
	sendPrompt(sender)
	log.Println(fmt.Sprintf("\x1b[1;37m[%s]:%s\x1b[0m", sender.name, msg))
	sendPrompt(receiver)
}

func isAPrompt(str string) (string, bool) {
	if strings.HasPrefix(str, "/rename ") {
		return "rename", true
	} else if strings.TrimRight(str, " ") == "/users" {
		return "users", true
	}
	return "", false
}
