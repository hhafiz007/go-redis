package main

import (
	"bufio"
	"flag"
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
			return
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
			if len(*&configValues.dir) > 0 {
				fileContent, _ := os.ReadFile(fmt.Sprintf("%s/%s", configValues.dir, configValues.dbfilename))
				fmt.Println(fileContent)
				_ = unMarshalRdb(fileContent)
			}
			key := string(redisArguments[0].bytes)
			value := myMap[key].value
			currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)

			if myMap[key].isLimited == true && currentTimeMillis > myMap[key].time {

				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))

			} else {

				fmt.Printf("Sending get to cliegnt with key and value and time%s %s %d\n", key, value, currentTimeMillis)
				conn.Write([]byte(fmt.Sprintf("+%s\r\n", value)))
			}
		case "config":
			configType := string(redisArguments[1].bytes)
			configEncoded := configValues.encodeConfigValues(configType)
			fmt.Printf("Sendiing Config to client : %s\n", configEncoded)
			conn.Write([]byte(configEncoded))
		case "keys":
			if len(*&configValues.dir) > 0 {
				fileContent, _ := os.ReadFile(fmt.Sprintf("%s/%s", configValues.dir, configValues.dbfilename))
				fmt.Println(fileContent)
				_ = unMarshalRdb(fileContent)
			}

			fmt.Println(myMap)

			keysCommand := string(redisArguments[0].bytes)
			var response string

			if keysCommand == "*" {
				response = getRespKeyArray(myMap)
			} else {
				response = fmt.Sprintf("$%d\r\n%s\r\n", len(myMap[keysCommand].value), myMap[keysCommand].value)
			}

			fmt.Println("The reesponse will", response)
			conn.Write([]byte(response))
			return

			// Convert the []byte to a string and print it
			// fileContentStr := string(fileContent)
			// fmt.Println(fileContent)

		}

	}

}

func main() {

	dir := flag.String("dir", "", "The directory where RDB files are stored")
	dbfilename := flag.String("dbfilename", "", "The name of the RDB file")
	flag.Parse()
	fmt.Println("dir:", *dir, len(*dir))
	fmt.Println("dbfilename:", *dbfilename)
	configValues = initConfigValues(dir, dbfilename)

	switch {

	default:

		l, err := net.Listen("tcp", "0.0.0.0:6379")
		fmt.Println("Listening to connections", l)
		if err != nil {
			fmt.Println("Failed to bind to port 6379")
			os.Exit(1)
		}

		myMap = make(map[string]redisValue)

		// fmt.Println(rdbFileData)

		// Will keep on running a for loop for accepting mu
		for {
			conn, err := l.Accept()

			if err != nil {
				fmt.Println("Failed to binad to port 6379")
				os.Exit(1)
			}
			go handleConnection(conn)
		}

	}

}
