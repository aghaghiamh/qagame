package middleware

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/claims"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/aghaghiamh/gocast/QAGame/pkg/timestamp"
	"github.com/aghaghiamh/gocast/QAGame/service/presenceservice"
	"github.com/labstack/echo/v4"
)

func Presence(presenceSvc presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, cErr := claims.GetClaimsFromEchoContext(c)
			if cErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Please provide a valid auth token")
			}

			_, err := presenceSvc.Upsert(
				c.Request().Context(),
				dto.PresenceUpsertRequest{
					UserID:    claims.UserID,
					Timestamp: timestamp.Now(),
				},
			)
			if err != nil {
				return echo.NewHTTPError(richerr.ErrUnexpected, err.Error())
			}

			return next(c)
		}
	}
}
