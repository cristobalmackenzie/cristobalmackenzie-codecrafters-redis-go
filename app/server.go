package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read the array length (we expect *1 for a simple PING)
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}

		// Read the length of the command
		_, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}

		// Read the actual command
		_, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}

		// assuming its a ping
		conn.Write([]byte("+PONG\r\n"))
	}
}
