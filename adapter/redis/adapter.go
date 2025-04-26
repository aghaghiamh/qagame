package redis

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RedisClient struct {
	driver *redis.Client
}

func New(redisConfig Config) RedisClient {
	return RedisClient{
		driver: redis.NewClient(&redis.Options{
			Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
		}),
	}
}

func (r RedisClient) Driver() *redis.Client {
	return r.driver
}
