package functions

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	
	name, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	fmt.Println(name, "has joined!")

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(name, "disconnected")
			return
		}
		fmt.Print("[", name, "]: ", msg)
	}
}
