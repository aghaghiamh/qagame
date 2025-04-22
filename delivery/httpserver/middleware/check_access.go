package middleware

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/claims"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/labstack/echo/v4"
)

func CheckAccess(authorizationSvc authorizationservice.Service, APIPermissionTitles ...entity.PermissionTitle) echo.MiddlewareFunc {
	// TODO: many other endpoints are required for authorization flow including source_permission and
	// access_controls creation, or seed them in migrations,
	// creating the backoffice users in backofficeSvc using one admin which we manually defined in our db, ...
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, cErr := claims.GetClaimsFromEchoContext(c)
			if cErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Please provide a valid auth token")
			}

			isPermitted, pErr := authorizationSvc.CheckPermissions(
				claims.UserID,
				claims.UserRole,
				APIPermissionTitles...,
			)
			if pErr != nil || !isPermitted {
				if pErr != nil && !isPermitted {
					code, msg := httpmapper.MapResponseCustomErrorToHttp(pErr)
					return echo.NewHTTPError(code, msg)
				}

				return echo.NewHTTPError(richerr.ErrUnauthorized, "You are not permitted to access this resource!")
			}

			return nil
		}
	}
}
