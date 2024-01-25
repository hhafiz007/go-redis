package main

import (
	"fmt"

	"net"
	"os"
)

func main() {

	serverAlive := []byte("+PONG\r\n")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	fmt.Println("Listening to connections", l)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for i := 0; i < 2; i++ {

		conn, _ := l.Accept()
		_, err = conn.Write(serverAlive)

	}

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

}
