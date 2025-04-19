package httpserver

import (
	"log"

	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Server struct {
	cfg         HttpConfig
	userHandler userhandler.UserHandler
}

func New(cfg HttpConfig, userHandler userhandler.UserHandler) Server {
	return Server{
		cfg:         cfg,
		userHandler: userHandler,
	}
}

func (s *Server) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s.userHandler.SetUserRoutes(e)

	if err := e.Start(s.cfg.Host + ":" + s.cfg.Port); err != nil {
		log.Fatalf("Couldn't Listen to the %s port.", s.cfg.Port)
	}
}
