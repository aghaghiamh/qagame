package main

import (
	"log"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/delivery/userserver"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	// user service registration
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
	userSvc := userservice.New(repo, &authSvc)
	handler := userserver.New(userSvc)

	e.POST("/user/register", handler.UserRegisterHandler)
	e.POST("/user/login", handler.UserLoginHandler)

	if err := e.Start("localhost:8080"); err != nil {
		log.Fatalf("Couldn't Listen to the 8080 port.")
	}
}
