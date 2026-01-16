package functions

import (
	"fmt"
	"strings"
	"time"
)

func saveHistory(msg string) {
	mu.Lock()
	defer mu.Unlock()
	history = append(history, msg)
}

func sendPrompt(c *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sendPromptLocked(c)
}

func sendPromptLocked(c *Client) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(c.conn, "\x1b[1;37m[%s][%s]:\x1b[0m", timestamp, c.name)
}

func validName(name string) bool {
	mu.Lock()
	defer mu.Unlock()
	if name == "" || len(name) > 10 {
		return false
	}
	for _, r := range name {
		// Si le caractère n'est PAS une minuscule ET n'est PAS une majuscule ET n'est PAS un chiffre
		if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
			return false
		}
	}
	for _, c := range clients {
		if c.name == name {
			return false
		}
	}
	return true
}

func validMsg(msg string) bool {
	mu.Lock()
	defer mu.Unlock()

	if msg == "" || len(msg) > 100 {
		return false
	}

	// Caractères autorisés : lettres, chiffres, espaces et ponctuation courante

	for _, r := range msg {
		if !strings.ContainsRune(allowedChars, r) {
			return false
		}
	}

	return true
}
