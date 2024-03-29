package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type config struct {
	dir        string
	dbfilename string
}

type rdbFile struct {
	magicString    string
	version        string
	dbSelector     int
	resizeDb       string
	AuxiliaryField string

	data map[string]redisValue
}

const (
	AuxiliaryField     int = 250
	ResizeDb           int = 251
	ExpiryMilliSeconds int = 252
	ExpirySeconds      int = 253
	DatabaseSelector   int = 254
	fileEnd            int = 255
	stringEncoding     int = 0
)

var configValues config

func initConfigValues(dir *string, dbfilename *string) config {

	return config{
		dir:        *dir,
		dbfilename: *dbfilename,
	}

}

func (cn *config) encodeConfigValues(configType string) string {

	var configTypeValue string

	if configType == "dir" {
		configTypeValue = cn.dir
	} else {
		configTypeValue = cn.dbfilename
	}

	return fmt.Sprintf("*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(configType), configType, len(configTypeValue), configTypeValue)
}

func unMarshalRdb(fileCont []byte) rdbFile {

	rdbDumpData := rdbFile{
		data: make(map[string]redisValue),
	}
	rdbDumpData.magicString = string(fileCont[:5])
	rdbDumpData.version = string(fileCont[5:9])

	fmt.Println("Welcome to unmarshal with magicString and version", rdbDumpData.magicString, rdbDumpData.version)

	currentInd := 9

	for {

		switch int(fileCont[currentInd]) {

		case AuxiliaryField:
			rdbDumpData.AuxiliaryField, currentInd = rdbDumpData.handleAuxilliary(currentInd, fileCont)
		case DatabaseSelector:
			rdbDumpData.dbSelector, currentInd = rdbDumpData.handledtabaseSelector(currentInd, fileCont)
		case ResizeDb:
			rdbDumpData.resizeDb, currentInd = rdbDumpData.handleResizeDb(currentInd, fileCont)
		case stringEncoding:
			currentInd = rdbDumpData.handleKeyValue(currentInd, fileCont, false, 0)
		case ExpiryMilliSeconds:
			currentInd = rdbDumpData.handleExpiryMiliSeconds(currentInd, fileCont)
		case ExpirySeconds:
			currentInd = rdbDumpData.handleExpirySeconds(currentInd, fileCont)

		default:
			break

		}

		if int(fileCont[currentInd]) == fileEnd {
			fmt.Println("breaking")
			break
		}

	}

	return rdbDumpData

}

func (rdb *rdbFile) handleAuxilliary(currentInd int, fileCont []byte) (string, int) {

	prevInd := currentInd + 1

	for {
		if int(fileCont[currentInd]) == DatabaseSelector {
			break
		}
		currentInd += 1
	}
	fmt.Printf("Auxiliary String is %s and next Index is %d ", string(fileCont[prevInd:currentInd]), currentInd)
	return string(fileCont[prevInd:currentInd]), currentInd

}

func (rdb *rdbFile) handledtabaseSelector(currentInd int, fileCont []byte) (int, int) {

	// prevInd := currentInd

	for {
		if int(fileCont[currentInd]) == ResizeDb {
			break
		}
		currentInd += 1
	}
	return 4, currentInd

}

func getLength(currentInd int, fileCont []byte) (int, int) {

	length := 0

	firstTwobits := (fileCont[currentInd] >> 6) & 0b11

	switch firstTwobits {

	case 0:
		length = int((fileCont[currentInd]) & 0b00111111)
		currentInd += 1
	case 1:
		byte1 := uint16(int((fileCont[currentInd])&0b00111111)) << 8
		byte2 := uint16(int(fileCont[currentInd+1]))

		length = int(byte1 + byte2)

		currentInd += 2
	case 2:

		length, _ = strconv.Atoi(string(fileCont[currentInd+1 : currentInd+5]))
		currentInd += 5

	}

	return length, currentInd

}

func (rdb *rdbFile) handleResizeDb(currentInd int, fileCont []byte) (string, int) {

	prevInd := currentInd
	currentInd += 1
	hashSize, currentInd := getLength(currentInd, fileCont)
	expirySize, currentInd := getLength(currentInd, fileCont)
	fmt.Println("The hash size and expiry size is  with last Index", hashSize, expirySize, currentInd)
	return string(fileCont[prevInd:currentInd]), currentInd

}

func (rdb *rdbFile) handleKeyValue(currentInd int, fileCont []byte, hasExpiry bool, time int64) int {

	currentInd += 1

	// fmt.Printf("index value is ,%d", int(fileCont[currentInd]))

	// keySizeInd := currentInd
	keyLength, currentInd := getLength(currentInd, fileCont)
	key := string(fileCont[currentInd : currentInd+keyLength])
	currentInd += keyLength
	valueLength, currentInd := getLength(currentInd, fileCont)
	value := string(fileCont[currentInd : currentInd+valueLength])

	// fmt.Println("The Keysize and valuesize from rdb file are", keySizeInd, valueSizeInd)

	rdb.data[key] = redisValue{
		value:     value,
		time:      time,
		isLimited: hasExpiry,
	}

	myMap[key] = redisValue{
		value:     value,
		time:      time,
		isLimited: hasExpiry,
	}
	fmt.Printf("The Key and value  and time from rdb file are %s and %s  and %d\n", key, value, time)

	currentInd += valueLength
	// fmt.Printf("index value is ,%d", int(fileCont[currentInd]))

	return currentInd

}
func (rdb *rdbFile) handleExpirySeconds(currentInd int, fileCont []byte) int {

	currentInd += 1
	expiryTime := int64(binary.LittleEndian.Uint32(fileCont[currentInd : currentInd+4]))
	nextInd := currentInd + 4
	expiryMiliSeconds := expiryTime * 1000
	fmt.Printf("The expiry in milliseconds from seconds is is %d\n", expiryTime)

	return rdb.handleKeyValue(nextInd, fileCont, true, int64(expiryMiliSeconds))

}
func (rdb *rdbFile) handleExpiryMiliSeconds(currentInd int, fileCont []byte) int {

	currentInd += 1

	expiryTime := int64(binary.LittleEndian.Uint64(fileCont[currentInd : currentInd+8]))

	nextInd := currentInd + 8

	fmt.Printf("The time and in miliSecons is %d\n", expiryTime)
	return rdb.handleKeyValue(nextInd, fileCont, true, int64(expiryTime))

}
