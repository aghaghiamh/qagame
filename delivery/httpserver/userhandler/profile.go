package userhandler

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/claims"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetProfileHandler(c echo.Context) error {
	claims, err := claims.GetClaimsFromEchoContext(c)
	if err != nil {
		echo.NewHTTPError(http.StatusUnauthorized, "Please provide a valid auth token")
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
