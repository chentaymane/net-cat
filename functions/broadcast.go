package functions

import (
	"fmt"
	"net"
)

// Broadcast a message to all clients except the sender
func broadcast(msg string, senderConn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range clients {

		// Skip sender
		if c.conn == senderConn {
			continue
		}

		// Send message to client
		fmt.Fprintln(c.conn, "\n"+msg)

		// Send prompt again (mutex already locked)
		sendPromptLocked(c)
	}
}
