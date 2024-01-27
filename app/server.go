package main

import (
	"bufio"
	"fmt"

	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		fmt.Println("Listening to connection", conn)
		redisMessage, err := handleRedisMessage(bufio.NewReader(conn))
		if err != nil {
			fmt.Println("Failed to bind to port 6379")
			os.Exit(1)
		}
		fmt.Println(redisMessage)
		redisCommand := string(redisMessage.array[0].bytes)
		fmt.Println(redisCommand)
		// redisArguments := redisMessage.Array()[1:]

		switch redisCommand {
		default:
			return
		case "ping":
			fmt.Printf("Sending ping to client\n")
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			fmt.Printf("Sending echo to client\n")
			conn.Write([]byte("+PONG\r\n"))

		}

	}

}

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	fmt.Println("Listening to connections", l)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	// Will keep on running a for loop for accepting mu
	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Failed to bind to port 6379")
			os.Exit(1)
		}
		go handleConnection(conn)
	}

}
