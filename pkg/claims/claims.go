package claims

import (
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/pkg/constant"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) (*authservice.Claims, error) {
	rawClaims := c.Get(constant.AuthMiddlewareContextKey)
	if claims, ok := rawClaims.(*authservice.Claims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("malwared jwt")
	}
}
