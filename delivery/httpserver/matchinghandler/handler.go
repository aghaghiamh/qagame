package matchinghandler

import (
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/matchingservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/matchingvalidator"
)

type Handler struct {
	matchingSvc      matchingservice.Service
	validator        matchingvalidator.MatchingValidator
	authSvc          authservice.Service
	authConfig       authservice.AuthConfig
	authorizationSvc authorizationservice.Service
}

func New(matchingSvc matchingservice.Service, authSvc authservice.Service,
	validator matchingvalidator.MatchingValidator, authConfig authservice.AuthConfig) Handler {
	return Handler{
		matchingSvc: matchingSvc,
		authSvc:     authSvc,
		validator:   validator,
		authConfig:  authConfig,
	}
}
