package backofficeuserhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) ListAllUsersHandler(c echo.Context) error {
	_, err := h.backofficeUserSvc.ListAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
