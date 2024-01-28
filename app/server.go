package main

import (
	"bufio"
	"fmt"
	"strconv"
	"time"

	"net"
	"os"
)

var myMap map[string]redisValue

type redisValue struct {
	value     string
	time      int64
	isLimited bool
}

func setRedisValue(val string, limit []byte) redisValue {
	currentTimeNano := time.Now().UnixNano()

	deadLine, _ := strconv.Atoi(string(limit))

	// Convert nanoseconds to milliseconds
	currentTimeMillis := (currentTimeNano / int64(time.Millisecond)) + int64(deadLine)

	isLimitedLocal := true
	if deadLine == 0 {
		isLimitedLocal = false
	}

	return redisValue{
		value:     val,
		time:      currentTimeMillis,
		isLimited: isLimitedLocal,
	}

}

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
			fmt.Printf("Sending echoo to client\n")
			// conn.Write(redisArguments[0].bytes)
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(string(redisArguments[0].bytes)), string(redisArguments[0].bytes))))
		case "set":

			key := string(redisArguments[0].bytes)
			value := string(redisArguments[1].bytes)
			arg3 := []byte("0")
			fmt.Println("Value of set is ", value)

			if len(redisArguments) >= 3 {
				arg3 = redisArguments[3].bytes
			}

			myMap[key] = setRedisValue(value, arg3)

			fmt.Printf("Sending set ato client with key and value and time %s %s %d\n", key, myMap[key].value, myMap[key].time)
			conn.Write([]byte("+OK\r\n"))
		case "get":
			key := string(redisArguments[0].bytes)
			value := myMap[key].value
			currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
			if myMap[key].isLimited == true && currentTimeMillis > myMap[key].time {

				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))

			} else {

				fmt.Printf("Sending get to cliegnt with key and value %s %s\n", key, value)
				conn.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
			}

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

	myMap = make(map[string]redisValue)

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
