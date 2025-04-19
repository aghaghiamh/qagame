package userhandler

import (
	"fmt"
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/constant"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"

	"github.com/labstack/echo/v4"
)

func getClaims(c echo.Context) (*authservice.Claims, error) {
	rawClaims := c.Get(constant.AUthMiddlewareSecretKey)
	if claims, ok := rawClaims.(*authservice.Claims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("malwared jwt")
	}
}

func (h UserHandler) GetProfileHandler(c echo.Context) error {
	claims, err := getClaims(c)
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
