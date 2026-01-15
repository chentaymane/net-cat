package functions

import (
	"fmt"
	"time"
)

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

	for _, c := range clients {
		if c.name == name {
			return false
		}
	}
	return true
}
