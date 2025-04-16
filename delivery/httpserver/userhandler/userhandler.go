package userhandler

import (
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	userservice "github.com/aghaghiamh/gocast/QAGame/service/userservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/uservalidator"
)

type UserHandler struct {
	userSvc   userservice.Service
	authSvc   authservice.Service
	validator uservalidator.UserValidator
}

func New(userSvc userservice.Service, authSvc authservice.Service, validator uservalidator.UserValidator) UserHandler {
	return UserHandler{
		userSvc:   userSvc,
		authSvc:   authSvc,
		validator: validator,
	}
}
