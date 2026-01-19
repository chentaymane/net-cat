package functions

import (
	"fmt"
	"net"
)
func broadcastToOthers(msg string, senderConn net.Conn) {
	mu.Lock()
	
	defer mu.Unlock()

	for _, c := range clients {
		if c.conn == senderConn {
			continue
		}

		

		if c.isRenaming {
			// buffer instead of losing
			c.pending = append(c.pending, msg)
			continue
		}

		fmt.Fprintln(c.conn, "\n"+msg)
		sendPromptLocked(c)

		
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
		mu.Lock()
		// Print newline first to clear the current prompt line
		fmt.Fprintln(c.conn, "\n"+msg)
		sendPromptLocked(c)
		mu.Unlock()
	}
}
