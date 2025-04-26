package main

import (
	redisAdapter "github.com/aghaghiamh/gocast/QAGame/adapter/redis"
	"github.com/aghaghiamh/gocast/QAGame/config"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/backofficeuserhandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/matchinghandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"

	// "github.com/aghaghiamh/gocast/QAGame/repository/migrator"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	accesscontroldb "github.com/aghaghiamh/gocast/QAGame/repository/mysql/accesscontrol"
	userdb "github.com/aghaghiamh/gocast/QAGame/repository/mysql/user"
	matchingdb "github.com/aghaghiamh/gocast/QAGame/repository/redis/matching"
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/backofficeuserservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/matchingvalidator"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"

	"github.com/aghaghiamh/gocast/QAGame/service/matchingservice"
)

func main() {
	config := config.LoadConfig()

	// General DB Connector
	generalMysqlDB, _ := mysql.New(config.DBConfig)

	// Redis Adapter
	redisAdapter := redisAdapter.New(config.RedisConfig)

	// Auth Service
	authSvc := authservice.New(config.AuthConfig)

	// Access-Control / Authorization Service
	acRepo := accesscontroldb.New(generalMysqlDB)
	authorizationSvc := authorizationservice.New(acRepo)

	// User Service
	userRepo := userdb.New(generalMysqlDB)
	uservalidator := uservalidator.New(userRepo)
	userSvc := userservice.New(userRepo, &authSvc)
	userHandler := userhandler.New(userSvc, authSvc, uservalidator, config.AuthConfig)

	// Matching Service
	matchingRepo := matchingdb.New(redisAdapter)
	matchingSvc := matchingservice.New(matchingRepo, config.MatchingConfig)
	matchingValidator := matchingvalidator.New(matchingRepo)
	matchingHandler := matchinghandler.New(matchingSvc, authSvc, matchingValidator, config.AuthConfig)

	// Back-Office Service
	backofficeUserSvc := backofficeuserservice.New()
	backofficeUserHandler := backofficeuserhandler.New(backofficeUserSvc, authSvc, config.AuthConfig, authorizationSvc)

	// m := migrator.New("mysql", config.DBConfig)
	// m.Up()

	server := httpserver.New(config.ServerConfig, userHandler, backofficeUserHandler, matchingHandler)
	server.Serve()
}
