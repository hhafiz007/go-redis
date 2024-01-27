package main

import (
	"bufio"
	"fmt"
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

	dataType, err := stream.ReadByte()
	fmt.Println("Data Type", dataType)
	if err != nil {
		fmt.Println("Error while reading dataType")
		return RedisMessage{}, err
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
		fmt.Println("Error reading array")
		return RedisMessage{}, err
	}

	items := int(bytes[0])

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

	data := make([]byte, int(bytes[0])+2)

	_, err = stream.Read(data)

	if err != nil {
		fmt.Println("Error reading bulk string")
		return RedisMessage{}, err
	}

	myMsg := RedisMessage{
		typ:   BulkString,
		bytes: data[:len(data)-2],
	}

	fmt.Println("Bulk String is", myMsg)

	return myMsg, nil

}
