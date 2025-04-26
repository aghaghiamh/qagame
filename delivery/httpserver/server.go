package httpserver

import (
	"log"

	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/backofficeuserhandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/matchinghandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Server struct {
	cfg               HttpConfig
	userHandler       userhandler.Handler
	backofficeHandler backofficeuserhandler.Handler
	matchingHandler   matchinghandler.Handler
}

func New(cfg HttpConfig, userHandler userhandler.Handler, backofficeHandler backofficeuserhandler.Handler,
	matchingHandler matchinghandler.Handler) Server {
	return Server{
		cfg:               cfg,
		userHandler:       userHandler,
		backofficeHandler: backofficeHandler,
		matchingHandler:   matchingHandler,
	}
}

func (s *Server) Serve() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s.userHandler.SetRoutes(e)
	s.backofficeHandler.SetRoutes(e)
	s.matchingHandler.SetRoutes(e)

	if err := e.Start(s.cfg.Host + ":" + s.cfg.Port); err != nil {
		log.Fatalf("Couldn't Listen to the %s port.", s.cfg.Port)
	}
}
