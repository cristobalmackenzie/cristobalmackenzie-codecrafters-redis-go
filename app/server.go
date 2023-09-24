package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue // instead of exiting, just handle the next connection
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		command, params, err := readCommand(reader)

		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			return
		}
		if command == "ping" {
			conn.Write([]byte("+PONG\r\n"))
		} else if command == "echo" {
			output := fmt.Sprintf("+%s\r\n", params[0])
			conn.Write([]byte(output))
		}
	}
}

func readCommand(reader *bufio.Reader) (string, [](string), error) {
	// *2, number of parameters
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", [](string){}, err
	}
	numParams := parseNumParams(line)

	params := [](string){}
	for i := 0; i < numParams; i++ {
		// $4, length of param
		_, err = reader.ReadString('\n')
		if err != nil {
			return "", [](string){}, err
		}

		// param itself
		line, err = reader.ReadString('\n')
		if err != nil {
			return "", [](string){}, err
		}
		params = append(params, strings.TrimSuffix(line, "\r\n"))
	}
	return params[0], params[1:], nil
}


func parseNumParams(inputStr string) int {
	line := strings.TrimPrefix(inputStr, "*")
	line = strings.TrimSuffix(line, "\r\n")
	numParams, err := strconv.Atoi(line)
	if err != nil {
		panic("Couldn't parse number of parameters")
	}
	return numParams
}
