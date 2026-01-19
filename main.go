package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"main/functions"
)

func main() {
	port := ""
	if len(os.Args) == 2 && os.Args[1] != "" {
		port = os.Args[1]
	} else {
		port = "8989"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("\x1b[38;5;198mServer listening on port " + port + "\x1b[0m")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go functions.HandleClient(conn)
	}
}
