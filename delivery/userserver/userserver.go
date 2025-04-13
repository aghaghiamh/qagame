package userserver

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (server *UserServer) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userGroup := e.Group("/user")
	userGroup.POST("/register", server.UserRegisterHandler)
	userGroup.POST("/login", server.UserLoginHandler)
	userGroup.GET("/profile", server.UserGetProfileHandler)

	if err := e.Start("localhost:8080"); err != nil {
		log.Fatalf("Couldn't Listen to the 8080 port.")
	}
}
