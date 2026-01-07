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
	
	// Ask for name
	fmt.Print("Enter your name: ")
	scanner.Scan() // MUST call Scan before Text
	name := scanner.Text()

	// Start live prompt loop
	for {
		fmt.Print("[", conn.LocalAddr(), "][", name, "]: ") // print prompt
		if !scanner.Scan() { // read user input
			break
		}
		text := scanner.Text()

		if text == "exit" {
			fmt.Println("Exiting...")
			return
		}

		// send name and message to server
		conn.Write([]byte(name + "\n"))
		conn.Write([]byte(text + "\n"))
	}
}
