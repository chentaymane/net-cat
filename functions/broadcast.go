package functions

import (
	"fmt"
	"net"
)

func broadcast(msg string, senderConn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range clients {
		if c.conn == senderConn {
			continue
		}
		fmt.Fprintln(c.conn, "\n"+msg)
		sendPromptLocked(c)
	}
}
