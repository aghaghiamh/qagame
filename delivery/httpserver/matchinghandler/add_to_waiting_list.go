package matchinghandler

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/claims"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/labstack/echo/v4"
)

func (h Handler) AddToWaitingListHandler(c echo.Context) error {
	claims, err := claims.GetClaimsFromEchoContext(c)
	if err != nil {
		echo.NewHTTPError(http.StatusUnauthorized, "Please provide a valid auth token")
	}

	req := dto.AddToWaitingListRequest{
		UserID: claims.UserID,
	}

	if err := c.Bind(&req); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fieldErrs, vErr := h.validator.ValidateAddToWaitingList(req)
	if vErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": errmsg.ErrMsgInvalidInput,
			"errors":  fieldErrs,
		})
	}

	resp, err := h.matchingSvc.AddToWaitingList(req)
	if err != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(err)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
