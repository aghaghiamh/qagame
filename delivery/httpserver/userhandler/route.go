package userhandler

import "github.com/labstack/echo/v4"

func (h UserHandler) SetUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/user")

	userGroup.POST("/register", h.RegisterHandler)
	userGroup.POST("/login", h.LoginHandler)
	userGroup.GET("/profile", h.GetProfileHandler)
}
