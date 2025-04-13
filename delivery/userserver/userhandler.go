package userserver

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	userservice "github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

	resp, regErr := h.userservice.Register(regReq)
	if regErr != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(regErr)

		return echo.NewHTTPError(code, msg)
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
		code, msg := httpmapper.MapResponseCustomErrorToHttp(err)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, loginResp)
}

func (h UserServer) UserGetProfileHandler(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")
	claims, vErr := h.authservice.VerifyToken(bearerToken)
	if vErr != nil {
		// TODO: Use Refresh Token
		if vErr == jwt.ErrTokenExpired {
		}
		code, msg := httpmapper.MapResponseCustomErrorToHttp(vErr)

		return echo.NewHTTPError(code, msg)
	}

	profileReq := userservice.UserProfileRequest{
		UserID: claims.UserID,
	}

	profileResp, respErr := h.userservice.GetProfile(profileReq)
	if respErr != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(respErr)
		
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, profileResp)
}
