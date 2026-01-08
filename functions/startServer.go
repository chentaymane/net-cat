package functions

import (
	"fmt"
	"net"
)

// SERVER

func StartServer(port string) {
	listener, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("New connection:", conn.RemoteAddr())

		
		go HandleClient(conn)
	}
}
