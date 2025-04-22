package userhandler

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/labstack/echo/v4"
)

func (h Handler) LoginHandler(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fieldErrs, vErr := h.validator.ValidateLoginRequest(req)
	if vErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": errmsg.ErrMsgInvalidInput,
			"errors":  fieldErrs,
		})
	}

	resp, err := h.userSvc.Login(req)
	if err != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(err)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
