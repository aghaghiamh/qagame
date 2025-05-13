package userhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "everything is ok")
}