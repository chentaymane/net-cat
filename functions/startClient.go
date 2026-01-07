package functions

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func StartClient(port string) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Connected to server!")
	fmt.Print("enter your name :")
	scanner.Scan()
	name := scanner.Text()
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			return
		}
		fmt.Print("[", conn.LocalAddr(), " ", name, "]: ", text)
		conn.Write([]byte(name + "\n"))
		conn.Write([]byte(text + "\n")) // ⬅ SEND TO SERVER
	}

	fmt.Println("Client exiting.")
}
