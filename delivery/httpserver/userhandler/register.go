package userhandler

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/labstack/echo/v4"
)

func (h UserHandler) RegisterHandler(c echo.Context) error {
	var regReq dto.RegisterRequest
	if err := c.Bind(&regReq); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fieldErrs, vErr := h.validator.ValidateRegisterRequest(regReq)
	if vErr != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": errmsg.ErrMsgInvalidInput,
			"errors":  fieldErrs,
		})
	}

	resp, regErr := h.userSvc.Register(regReq)
	if regErr != nil {
		code, msg := httpmapper.MapResponseCustomErrorToHttp(regErr)

		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
