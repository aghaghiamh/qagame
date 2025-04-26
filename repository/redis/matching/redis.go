package matching

import (
	redisAdapter "github.com/aghaghiamh/gocast/QAGame/adapter/redis"
)

type Storage struct {
	adapter redisAdapter.RedisClient
}

func New(adapter redisAdapter.RedisClient) *Storage {
	return &Storage{
		adapter: adapter,
	}
}
