package main

import (
	"fmt"

	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	serverAlive := []byte("+PONG\r\n")

	for {
		buf := make([]byte, 256)
		_, _ = conn.Read(buf)
		if len(buf) > 0 {
			conn.Write(serverAlive)
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

	for {
		conn, _ := l.Accept()
		defer conn.Close()
		go handleConnection(conn)
	}

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

}
