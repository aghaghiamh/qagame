package httpserver

import (
	"log"

	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	userHandler userhandler.UserHandler
}

func New(userHandler userhandler.UserHandler) Server {
	return Server{
		userHandler: userHandler,
	}
}

func (s *Server) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s.userHandler.SetUserRoutes(e)

	if err := e.Start("localhost:8080"); err != nil {
		log.Fatalf("Couldn't Listen to the 8080 port.")
	}
}
