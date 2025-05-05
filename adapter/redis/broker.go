package redis

import (
	"context"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/labstack/gommon/log"
)

func (c RedisClient) Publish(event entity.Event, payload string){
	const op = "redisAdapter.Publish"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	if err := c.driver.Publish(ctx, string(event), payload).Err(); err != nil {
		log.Errorf("%s: %v\n", op, err)
	}
}