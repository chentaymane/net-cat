package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"main/functions" // import your functions package
)

func main() {
	// Get port from arguments, default to 8989
	port := "8989"
	if len(os.Args) == 2 && os.Args[1] != "" {
		port = os.Args[1]
	}

	// Listen on TCP port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("\x1b[38;5;198mServer listening on port " + port + "\x1b[0m")

	for {
		// Accept new client connections
		conn, err := listener.Accept()
		if err != nil {
			continue // ignore failed connections
		}

		// Handle client in a new goroutine
		go functions.HandleClient(conn)
	}
}
