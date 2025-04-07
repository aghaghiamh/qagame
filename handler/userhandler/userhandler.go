package userhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	userservice "github.com/aghaghiamh/gocast/QAGame/service/userservice"
)

type UserHandler struct {
	userservice userservice.Service
}

func New(userservice userservice.Service) UserHandler {
	return UserHandler{userservice: userservice}
}

// TODO: Duplicated repo and service initialization
func (h UserHandler) UserRegisterHandler(wr http.ResponseWriter, req *http.Request) {
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

	_, err := h.userservice.Register(regReq)
	if err != nil {
		fmt.Print(err)
		return
	}

	// return response
	fmt.Fprint(wr, "successful registration")
}

func (h UserHandler) UserLoginHandler(wr http.ResponseWriter, req *http.Request) {
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

	loginResp, err := h.userservice.Login(loginReq)
	if err != nil {
		fmt.Print(err)
		return
	}

	// return response
	fmt.Fprintf(wr, "successful login, at: %s, rt: %s \n", loginResp.AccessToken, loginResp.RefreshToken)
}
