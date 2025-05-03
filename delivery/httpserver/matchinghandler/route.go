package matchinghandler

import (
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/matching")

	userGroup.POST("/add-to-waiting-list", h.AddToWaitingListHandler, middleware.Auth(h.authSvc, h.authConfig), middleware.Presence(h.presenceSvc))
	userGroup.GET("/get-list", h.GetFromWaitingListHandler)
}
