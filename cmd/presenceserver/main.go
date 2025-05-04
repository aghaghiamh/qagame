package main

import (
	redisAdapter "github.com/aghaghiamh/gocast/QAGame/adapter/redis"
	"github.com/aghaghiamh/gocast/QAGame/config"
	"github.com/aghaghiamh/gocast/QAGame/delivery/grpc/presenceserver"
	"github.com/aghaghiamh/gocast/QAGame/repository/redis/presence"
	"github.com/aghaghiamh/gocast/QAGame/service/presenceservice"
)

func main() {
	config := config.LoadConfig()

	// Redis Adapter
	redisAdapter := redisAdapter.New(config.Redis)

	// Presence Service
	presenceRepo := presence.New(redisAdapter)
	presenceSvc := presenceservice.New(config.PresenceSvc, presenceRepo)

	// Presence gRPC Server
	server := presenceserver.New(&presenceSvc)
	server.Serve()
}
