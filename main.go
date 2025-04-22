package main

import (
	"github.com/aghaghiamh/gocast/QAGame/config"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/backofficeuserhandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"

	// "github.com/aghaghiamh/gocast/QAGame/repository/migrator"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	accesscontroldb "github.com/aghaghiamh/gocast/QAGame/repository/mysql/accesscontrol"
	userdb "github.com/aghaghiamh/gocast/QAGame/repository/mysql/user"
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/backofficeuserservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"
)

func main() {
	config := config.LoadConfig()

	generalMysqlDB, _ := mysql.New(config.DBConfig)
	userRepo := userdb.New(generalMysqlDB)

	authSvc := authservice.New(config.AuthConfig)

	uservalidator := uservalidator.New(userRepo)

	acRepo := accesscontroldb.New(generalMysqlDB)
	authorizationSvc := authorizationservice.New(acRepo)

	userSvc := userservice.New(userRepo, &authSvc)
	userHandler := userhandler.New(userSvc, authSvc, uservalidator, config.AuthConfig)

	backofficeUserSvc := backofficeuserservice.New()
	backofficeUserHandler := backofficeuserhandler.New(backofficeUserSvc, authSvc, config.AuthConfig, authorizationSvc)

	// m := migrator.New("mysql", config.DBConfig)
	// m.Up()

	server := httpserver.New(config.ServerConfig, userHandler, backofficeUserHandler)
	server.Serve()
}
