package main

import (
	"bufio"
	"bytes"
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
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	// defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read the array length (we expect *1 for a simple PING)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}

		if line != "*1\r\n" {
			continue
		}

		// Read the length of the command
		line, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}

		if !bytes.HasPrefix([]byte(line), []byte("$")) {
			continue
		}

		// Read the actual command
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}

		if cmd == "PING\r\n" {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	reader := bufio.NewReader(conn)
// 	message, err := reader.ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Error reading from connection:", err.Error())
// 		return
// 	}

// 	fmt.Println("Received:", message)

// 	conn.Write([]byte("+PONG\r\n"))
// }