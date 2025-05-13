package userhandler

import (
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/user")

	userGroup.GET("/health", h.HealthCheckHandler)
	userGroup.POST("/register", h.RegisterHandler)
	userGroup.POST("/login", h.LoginHandler)
	userGroup.GET("/profile", h.GetProfileHandler, middleware.Auth(h.authSvc, h.authConfig))
}
