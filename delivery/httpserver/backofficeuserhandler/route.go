package backofficeuserhandler

import (
	"github.com/aghaghiamh/gocast/QAGame/delivery/httpserver/middleware"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/back-office")

	userGroup.GET("/users", h.ListAllUsersHandler,
		middleware.Auth(h.authSvc, h.authConfig),
		middleware.CheckAccess(h.authorizationSvc, entity.UserListPermission),
	)
}
