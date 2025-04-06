package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
		fmt.Fprint(wr, "Bad Request: only post request accepted for login")
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
	loginResp, err := service.Login(loginReq)
	if err != nil {
		fmt.Print(err)
		return
	}

	// return response
	fmt.Fprintf(wr, "successful login, token: %s", loginResp.Token)
}

func UserAuthHandler(wr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprint(wr, "Bad Request: only get request accepted for Authorization")
		return
	}

	// unmarshal request
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Print("Authorization header required", http.StatusUnauthorized)
		return
	}

	// Check if the header has the right format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		fmt.Print("Authorization header must be in the format: Bearer {token}", http.StatusUnauthorized)
		return
	}

	// Create auth request with the token
	authReq := userservice.AuthRequest{
		Token: parts[1],
	}

	service := userservice.New(nil)
	_, err := service.Authorize(authReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return response
	fmt.Fprint(wr, "successful auth")
}
