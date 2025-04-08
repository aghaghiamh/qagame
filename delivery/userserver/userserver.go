package userserver

import (
	"log"
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	userservice "github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserServer struct {
	userservice userservice.Service
	authservice authservice.Service
}

func New(userservice userservice.Service, authservice authservice.Service) UserServer {
	return UserServer{
		userservice: userservice,
		authservice: authservice,
	}
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

func (h UserServer) UserGetProfileHandler(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")
	claims, vErr := h.authservice.VerifyToken(bearerToken)
	if vErr != nil {
		// TODO: Use Refresh Token
		if vErr == jwt.ErrTokenExpired {}
		return vErr
	}

	profileReq := userservice.UserProfileRequest{
		UserID: claims.UserID,
	}

	profileResp, respErr := h.userservice.GetUserProfile(profileReq)
	if respErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, respErr.Error())
	}

	return c.JSON(http.StatusOK, profileResp)
}

func (server *UserServer) Serve(){
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userGroup := e.Group("/user")
	userGroup.POST("/register", server.UserRegisterHandler)
	userGroup.POST("/login", server.UserLoginHandler)
	userGroup.GET("/profile", server.UserGetProfileHandler)

	if err := e.Start("localhost:8080"); err != nil {
		log.Fatalf("Couldn't Listen to the 8080 port.")
	}
}