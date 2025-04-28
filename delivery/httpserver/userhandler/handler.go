package userhandler

import (
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	userservice "github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"
)

type Handler struct {
	userSvc    userservice.Service
	validator  uservalidator.UserValidator
	authSvc    authservice.Service
	authConfig authservice.Config
}

func New(userSvc userservice.Service, authSvc authservice.Service,
	validator uservalidator.UserValidator, authConfig authservice.Config) Handler {
	return Handler{
		userSvc:    userSvc,
		authSvc:    authSvc,
		validator:  validator,
		authConfig: authConfig,
	}
}
