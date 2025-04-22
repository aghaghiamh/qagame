package backofficeuserhandler

import (
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/backofficeuserservice"
)

type Handler struct {
	backofficeUserSvc backofficeuserservice.Service
	authSvc           authservice.Service
	authConfig        authservice.AuthConfig
	authorizationSvc  authorizationservice.Service
}

func New(
	backofficeUserSvc backofficeuserservice.Service,
	authSvc authservice.Service,
	authConfig authservice.AuthConfig,
	authorizationSvc authorizationservice.Service,
) Handler {
	return Handler{
		backofficeUserSvc: backofficeUserSvc,
		authSvc:           authSvc,
		authConfig:        authConfig,
		authorizationSvc:  authorizationSvc,
	}
}
