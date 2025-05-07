package httpserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/backofficeuserhandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/matchinghandler"
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/userhandler"
	"github.com/aghaghiamh/gocast/QAGame/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type HttpConfig struct {
	Host                    string        `mapstructure:"host"`
	Port                    string        `mapstructure:"port"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}

type Server struct {
	cfg               HttpConfig
	router            *echo.Echo
	userHandler       userhandler.Handler
	backofficeHandler backofficeuserhandler.Handler
	matchingHandler   matchinghandler.Handler
}

func New(cfg HttpConfig, userHandler userhandler.Handler, backofficeHandler backofficeuserhandler.Handler,
	matchingHandler matchinghandler.Handler) Server {
	return Server{
		cfg:               cfg,
		router:            echo.New(),
		userHandler:       userHandler,
		backofficeHandler: backofficeHandler,
		matchingHandler:   matchingHandler,
	}
}

func (s *Server) Serve() {
	s.router.Use(middleware.RequestID())
	
	s.router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogLatency: true,
		LogMethod: true,
		LogURI: true,
		LogRemoteIP: true,
		LogError: true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			var errMsg string
			if v.Error != nil {
				errMsg = v.Error.Error()
			}
			
			logger.Logger.Info("request",
				zap.String("request_id", v.RequestID),
				zap.Duration("latency", v.Latency),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.String("remote_ip", v.RemoteIP),
				zap.Error(v.Error),
				zap.String("error-msg", errMsg),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	s.router.Use(middleware.Recover())

	s.userHandler.SetRoutes(s.router)
	s.backofficeHandler.SetRoutes(s.router)
	s.matchingHandler.SetRoutes(s.router)

	if err := s.router.Start(s.cfg.Host + ":" + s.cfg.Port); err != nil {
		log.Fatalf("Couldn't Listen to the %s port.", s.cfg.Port)
	}
}

func (s Server) Shutdown() {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.router.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("error while shutting down the server: ", err)
	}

	fmt.Println("Gracefully shutdowned!!")
	<-ctxWithTimeout.Done()
}
