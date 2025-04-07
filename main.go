package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/handler/userhandler"
	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/userservice"
)

func main() {
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
	handler := userhandler.New(userSvc)

	http.HandleFunc("/user/register", handler.UserRegisterHandler)
	http.HandleFunc("/user/login", handler.UserLoginHandler)

	fmt.Print("Server is running on port 8080...")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatalf("Couldn't Listen to the 8080 port.")
	}
}
