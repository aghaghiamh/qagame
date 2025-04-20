package main

import (
	"github.com/aghaghiamh/gocast/QAGame/config"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	// "github.com/aghaghiamh/gocast/QAGame/repository/migrator"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"
)

func main() {
	config := config.LoadConfig()

	repo, _ := mysql.New(config.DBConfig)

	authSvc := authservice.New(config.AuthConfig)

	uservalidator := uservalidator.New(repo)
	userSvc := userservice.New(repo, &authSvc)
	userHandler := userhandler.New(userSvc, authSvc, uservalidator, config.AuthConfig)

	// m := migrator.New("mysql", config.DBConfig)
	// m.Up()

	server := httpserver.New(config.ServerConfig, userHandler)
	server.Serve()
}
