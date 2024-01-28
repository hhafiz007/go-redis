package main

import (
	"bufio"
	"fmt"

	"net"
	"os"
)

var myMap map[string]string

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		fmt.Println("Listening to connection", conn)
		redisMessage, err := handleRedisMessage(bufio.NewReader(conn))
		if err != nil {
			fmt.Println("Failed to bind to port 6379")
			os.Exit(1)
		}
		fmt.Println("redis message is", redisMessage)
		redisCommand := string(redisMessage.array[0].bytes)
		fmt.Println("redis command is", redisCommand)
		redisArguments := redisMessage.array[1:]

		switch redisCommand {
		default:
			fmt.Printf("Sending anwer to client\n")
			conn.Write([]byte(fmt.Sprintf("-ERR INVALID COMMAND %s\r\n", redisCommand)))
		case "ping":
			fmt.Printf("Sending ping to client\n")
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			fmt.Printf("Sending echo to client\n")
			// conn.Write(redisArguments[0].bytes)
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(string(redisArguments[0].bytes)), string(redisArguments[0].bytes))))
		case "set":

			key := string(redisArguments[0].bytes)
			value := string(redisArguments[1].bytes)
			myMap[key] = value
			fmt.Printf("Sending set to client with key and value %s %s\n", key, value)
			conn.Write([]byte("+OK\r\n"))

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

	myMap = make(map[string]string)

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
