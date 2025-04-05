package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/repository/mysql"
	userservice "github.com/aghaghiamh/gocast/QAGame/service"
)

// TODO: Duplicated repo and service initialization
func UserRegisterHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprint(wr, "Bad Request: only post request accepted for registration")
		return
	}

	// unmarshal request
	var regReq userservice.RegisterRequest
	data, rErr := io.ReadAll(req.Body)
	if rErr != nil {
		fmt.Fprintf(wr, "Unable to read body of request: %s", req.Body)
	}

	if err := json.Unmarshal(data, &regReq); err != nil {
		fmt.Fprintf(wr, "Bad Request: %s", data)
		return
	}

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
	service := userservice.New(repo)
	_, err := service.Register(regReq)
	if err != nil {
		fmt.Print(err)
		return
	}

	// return response
	fmt.Fprint(wr, "successful registration")
}

func UserLoginHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Fprint(wr, "Bad Request: only post request accepted for registration")
		return
	}

	// unmarshal request
	var loginReq userservice.LoginRequest
	data, rErr := io.ReadAll(req.Body)
	if rErr != nil {
		fmt.Fprintf(wr, "Unable to read body of request: %s", req.Body)
	}

	if err := json.Unmarshal(data, &loginReq); err != nil {
		fmt.Fprintf(wr, "Bad Request: %s", data)
		return
	}

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
	service := userservice.New(repo)
	_, err := service.Login(loginReq)
	if err != nil {
		fmt.Print(err)
		return
	}

	// return response
	fmt.Fprint(wr, "successful login")
}