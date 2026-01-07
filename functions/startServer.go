package functions

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// SERVER

func StartServer(port string) {
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening on port", port)

	// Accept one connection
	conn, err := listener.Accept()
	fmt.Println("New connection !")
	if err != nil {
		panic(err)
	}
 firstTime := true
	for {
		reader := bufio.NewReader(conn)
		name, err := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if  firstTime {
			fmt.Println(name, "Has joined !")
			firstTime =false
		}
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
		msgs, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Print("[", conn.LocalAddr(), "][", name, "]: ", msgs)
		}
	}

