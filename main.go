package main

import (
	"log"
	"net"

	"main/functions"
)

func main() {
	listener, err := net.Listen("tcp", ":8989")
	if err != nil {
		panic(err)
	}
	log.Println("Server listening on port 8989")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go functions.HandleClient(conn)
	}
}
