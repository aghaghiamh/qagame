package userhandler

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h UserHandler) GetProfileHandler(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")
	claims, vErr := h.authSvc.VerifyToken(bearerToken)
	if vErr != nil {
		// TODO: Use Refresh Token
		if vErr == jwt.ErrTokenExpired {
		}
		code, msg := httpmapper.MapResponseCustomErrorToHttp(vErr)

		return echo.NewHTTPError(code, msg)
	}

	profileReq := dto.UserProfileRequest{
		UserID: claims.UserID,
	}

	profileResp, respErr := h.userSvc.GetProfile(profileReq)
	if respErr != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(respErr)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, profileResp)
}
