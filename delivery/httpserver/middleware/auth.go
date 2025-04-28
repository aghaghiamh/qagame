package middleware

import (
	"github.com/aghaghiamh/gocast/QAGame/pkg/constant"
	"github.com/aghaghiamh/gocast/QAGame/pkg/httpmapper"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/golang-jwt/jwt/v5"
	ejwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(authSvc authservice.Service, authConf authservice.Config) echo.MiddlewareFunc {
	return ejwt.WithConfig(ejwt.Config{
		ContextKey: constant.AuthMiddlewareContextKey,
		SigningKey: authConf.SignKey,
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, vErr := authSvc.VerifyToken(auth)
			if vErr != nil {
				// TODO: Use Refresh Token
				if vErr == jwt.ErrTokenExpired {
				}
				code, msg := httpmapper.MapResponseCustomErrorToHttp(vErr)

				return nil, echo.NewHTTPError(code, msg)
			}
			return claims, nil
		},
	},
	)
}
