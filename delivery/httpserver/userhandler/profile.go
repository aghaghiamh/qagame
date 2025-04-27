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

	// TODO: here is an implementation of context without timeout or further functionalities, should be implemented in other handlers, too.
	profileResp, respErr := h.userSvc.GetProfile(c.Request().Context(), profileReq)
	if respErr != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(respErr)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, profileResp)
}
