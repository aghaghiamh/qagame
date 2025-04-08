package userserver

import (
	"net/http"

	userservice "github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/labstack/echo/v4"
)

type UserServer struct {
	userservice userservice.Service
}

func New(userservice userservice.Service) UserServer {
	return UserServer{userservice: userservice}
}

func (h UserServer) UserRegisterHandler(c echo.Context) error {
	var regReq userservice.RegisterRequest
	if err := c.Bind(&regReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.userservice.Register(regReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h UserServer) UserLoginHandler(c echo.Context) error {
	var loginReq userservice.LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	loginResp, err := h.userservice.Login(loginReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, loginResp)
}
