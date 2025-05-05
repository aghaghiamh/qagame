package main

import (
	"context"
	"fmt"
	"time"

	redisAdapter "github.com/aghaghiamh/gocast/QAGame/adapter/redis"
	"github.com/aghaghiamh/gocast/QAGame/config"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/eventencoder"
	"github.com/aghaghiamh/gocast/QAGame/pkg/typemapper"
	"github.com/labstack/gommon/log"
)

func main() {
	const op = "gamesubscriber"

	config := config.LoadConfig()

	// Redis Adapter
	redisAdapter := redisAdapter.New(config.Redis)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	fmt.Println("starting listening on channels")

	for {
		sub := redisAdapter.Driver().Subscribe(ctx, string(entity.MatchingMatchedUsersEvent))
		msg, rErr := sub.ReceiveMessage(ctx)
		if rErr != nil {
			log.Errorf("operation: %s, err: %v", op, rErr)
			return
		}

		switch msg.Channel {
		case string(entity.MatchingMatchedUsersEvent):
			processUsersMatchedEvent(msg.Payload)
		default:
			fmt.Println("invalid channel: ", msg.Channel)
		}
	}
}

func processUsersMatchedEvent(payloadStr string) {
	pbMp, dErr := eventencoder.MatchedPlayerUsersDecoder(payloadStr)
	if dErr != nil {
		log.Error(dErr)
	}

	mp := entity.MatchedPlayers{
		Category: entity.Category(pbMp.Category),
		UserIDs: typemapper.ArrayMapper(pbMp.UserIds, func(uID uint64) uint {
			return uint(uID)
		}),
	}

	fmt.Print("decoded payload: ", mp)
}
