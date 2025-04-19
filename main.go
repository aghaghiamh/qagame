package main

import (
	"time"

	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"
)

func main() {
	dbConf := mysql.MysqlConfig{
		Host:     "127.0.0.1",
		Port:     "3308",
		Username: "root",
		Password: "12345",
		DBName:   "users",

		MaxLifeTime: time.Minute * 3,
		MaxOpenConn: 10,
		MaxIdleConn: 10,
	}

	repo, _ := mysql.New(dbConf)

	authConf := authservice.AuthConfig{
		SignKey:              "secret-key",
		AccessSubject:        "at",
		RefreshSubject:       "rt",
		AccessTokenDuration:  time.Hour * 24,
		RefreshTokenDuration: time.Hour * 24 * 7,
	}
	authSvc := authservice.New(authConf)

	uservalidator := uservalidator.New(repo)

	userSvc := userservice.New(repo, &authSvc)
	userHandler := userhandler.New(userSvc, authSvc, uservalidator, authConf)

	server := httpserver.New(userHandler)
	server.Serve()
}
