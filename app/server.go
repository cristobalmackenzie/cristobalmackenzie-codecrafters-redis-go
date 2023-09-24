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
	listen()
}

func listen() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	rs := NewRedisStore()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue // instead of exiting, just handle the next connection
		}

		go handleConnection(conn, rs)
	}
}

func handleConnection(conn net.Conn, rs *RedisStore) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		command, params, err := readCommand(reader)

		if err != nil {
			// fmt.Println("Error reading from command:", err.Error())
			return
		}

		response := "+OK\r\n"
		if command == "ping" {
			response = "+PONG\r\n"
		} else if command == "echo" {
			response = fmt.Sprintf("+%s\r\n", params[0])
		} else if command == "set" {
			var px *int64
			if len(params) == 4 && params[2] == "px" {
				pxMillis, _ := strconv.Atoi(params[3])
				pxMillis64 := int64(pxMillis)
				px = &pxMillis64
			}
			rs.Set(params[0], params[1], px)
		} else if command == "get" {
			value, exists := rs.Get(params[0])
			if exists {
				response = fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
			} else {
				response = "$-1\r\n" 
			}
		}

		conn.Write([]byte(response))
	}
}

func readCommand(reader *bufio.Reader) (string, [](string), error) {
	// *2, number of parameters
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", [](string){}, err
	}
	numParams, err := parseNumParams(line)
	if err != nil {
		return "", [](string){}, err
	}

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


func parseNumParams(inputStr string) (int, error) {
	line := strings.TrimPrefix(inputStr, "*")
	line = strings.TrimSuffix(line, "\r\n")
	numParams, err := strconv.Atoi(line)
	if err != nil {
		return 0, err
	}
	return numParams, nil
}
