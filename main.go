package main

import (
	"os"
	"os/signal"

	redisAdapter "github.com/aghaghiamh/gocast/QAGame/adapter/redis"
	"github.com/aghaghiamh/gocast/QAGame/config"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/backofficeuserhandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/matchinghandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	"github.com/aghaghiamh/gocast/QAGame/scheduler"

	// "github.com/aghaghiamh/gocast/QAGame/repository/migrator"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	accesscontroldb "github.com/aghaghiamh/gocast/QAGame/repository/mysql/accesscontrol"
	userdb "github.com/aghaghiamh/gocast/QAGame/repository/mysql/user"
	matchingdb "github.com/aghaghiamh/gocast/QAGame/repository/redis/matching"
	"github.com/aghaghiamh/gocast/QAGame/repository/redis/presence"
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/backofficeuserservice"
	"github.com/aghaghiamh/gocast/QAGame/service/presenceservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/matchingvalidator"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"

	"github.com/aghaghiamh/gocast/QAGame/service/matchingservice"
)

func main() {
	config := config.LoadConfig()

	// General DB Connector
	generalMysqlDB, _ := mysql.New(config.DB)

	// Redis Adapter
	redisAdapter := redisAdapter.New(config.Redis)

	// m := migrator.New("mysql", config.DBConfig)
	// m.Up()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// run the http server
	userHandler, matchingSvc, matchingHandler, backofficeUserHandler := setup(config, generalMysqlDB, redisAdapter)
	server := httpserver.New(config.Server, userHandler, backofficeUserHandler, matchingHandler)
	go func() {
		server.Serve()
	}()

	// run the cronjob scheduler
	schDoneCH := make(chan bool)
	sch := scheduler.New(config.Scheduler, schDoneCH, matchingSvc)
	go func() {
		sch.Start()
	}()

	// Graceful Termination - wait until there is a os.signal on the quit channel then revoke all other children.
	<-quit
	schDoneCH <- true

	server.Shutdown()
}

func setup(config config.Config, mysqlDB *mysql.MysqlDB, redisAdapter redisAdapter.RedisClient) (
	userhandler.Handler, matchingservice.Service, matchinghandler.Handler, backofficeuserhandler.Handler) {
	// Auth Service
	authSvc := authservice.New(config.AuthSvc)

	// Access-Control / Authorization Service
	acRepo := accesscontroldb.New(mysqlDB)
	authorizationSvc := authorizationservice.New(acRepo)

	// Presence Service
	presenceRepo := presence.New(redisAdapter)
	presenceSvc := presenceservice.New(config.PresenceSvc, presenceRepo)

	// User Service
	userRepo := userdb.New(mysqlDB)
	uservalidator := uservalidator.New(userRepo)
	userSvc := userservice.New(userRepo, &authSvc)
	userHandler := userhandler.New(userSvc, authSvc, uservalidator, config.AuthSvc)

	// Matching Service
	matchingRepo := matchingdb.New(redisAdapter)
	matchingSvc := matchingservice.New(matchingRepo, config.MatchingSvc)
	matchingValidator := matchingvalidator.New(matchingRepo)
	matchingHandler := matchinghandler.New(matchingSvc, authSvc, presenceSvc, matchingValidator, config.AuthSvc)

	// Back-Office Service
	backofficeUserSvc := backofficeuserservice.New()
	backofficeUserHandler := backofficeuserhandler.New(backofficeUserSvc, authSvc, config.AuthSvc, authorizationSvc)

	return userHandler, matchingSvc, matchingHandler, backofficeUserHandler
}
