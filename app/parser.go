package main

import (
	"bufio"
	"fmt"
	"strconv"
)

type RedisMessageType byte

type RedisMessage struct {
	typ   RedisMessageType // type of message (SimpleString, BulkString, Array)
	bytes []byte
	array []RedisMessage
}

const (
	SimpleString RedisMessageType = '+'
	BulkString   RedisMessageType = '$'
	Array        RedisMessageType = '*'
)

func handleRedisMessage(stream *bufio.Reader) (RedisMessage, error) {
	fmt.Println("handling redis message")

	dataType, err := stream.ReadByte()

	for {
		fmt.Println("Data Type", dataType)
		if err != nil {
			fmt.Println("Error whiqle reading dataType", err)
			// return RedisMessage{}, err
		} else {
			break
		}

	}

	switch dataType {
	case byte(SimpleString):
		return parseSimpleString(stream)
	case byte(BulkString):
		return parseBulkString(stream)
	case byte(Array):
		return parseArray(stream)
	}

	return RedisMessage{}, nil

}

func parseArray(stream *bufio.Reader) (RedisMessage, error) {
	bytes, err := stream.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error readinig aarray", err)
		return RedisMessage{}, err
	}

	items, _ := strconv.Atoi(string(bytes[0]))

	myMsg := RedisMessage{
		typ: Array,
	}

	for i := 0; i < items; i += 1 {
		fmt.Println(i, items)
		curMsg, err := handleRedisMessage(stream)

		if err != nil {
			fmt.Println("Error reading array")
			return RedisMessage{}, err
		}
		myMsg.array = append(myMsg.array, curMsg)

	}

	fmt.Println("Array is", myMsg)

	return myMsg, nil

}

func parseSimpleString(stream *bufio.Reader) (RedisMessage, error) {

	currString, err := stream.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading SimpleString")
		return RedisMessage{}, err
	}

	myMsg := RedisMessage{
		typ:   SimpleString,
		bytes: currString[:len(currString)-2],
	}
	fmt.Println("Simple String is", myMsg)

	return myMsg, nil

}

func parseBulkString(stream *bufio.Reader) (RedisMessage, error) {
	bytes, err := stream.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error reading bulk string")
		return RedisMessage{}, err
	}

	length, _ := strconv.Atoi(string(bytes[0]))
	data := make([]byte, length+2)

	_, err = stream.Read(data)

	if err != nil {
		fmt.Println("Error reading bulk string")
		return RedisMessage{}, err
	}

	myMsg := RedisMessage{
		typ:   BulkString,
		bytes: data[:len(data)-2],
	}

	fmt.Println("Bulk String is", myMsg, string(myMsg.bytes))

	return myMsg, nil

}
