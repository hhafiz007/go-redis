package main

import "time"

type Storage struct {
	data map[string]RedisValue
}

type RedisValue struct {
	value      string
	expiration time.Time
}

func InitStorage() *Storage {
	return &Storage{
		data: make(map[string]RedisValue),
	}
}
