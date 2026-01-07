package main

import (
	"fmt"
	"os"

	"main/functions"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  Server: go run main.go server <port>")
		fmt.Println("  Client: go run main.go client <host:port>")
		return
	}

	mode := os.Args[1]
	address := os.Args[2]

	if mode == "server" {
		functions.StartServer(address)
	} else if mode == "client" {
		functions.StartClient(address)
	} else {
		fmt.Println("Unknown mode:", mode)
	}
}
