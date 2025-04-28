package matchinghandler

import (
	"github.com/aghaghiamh/gocast/QAGame/service/authorizationservice"
	"github.com/aghaghiamh/gocast/QAGame/service/authservice"
	"github.com/aghaghiamh/gocast/QAGame/service/matchingservice"
	"github.com/aghaghiamh/gocast/QAGame/service/presenceservice"
	"github.com/aghaghiamh/gocast/QAGame/validator/matchingvalidator"
)

type Handler struct {
	matchingSvc      matchingservice.Service
	authSvc          authservice.Service
	presenceSvc      presenceservice.Service
	validator        matchingvalidator.MatchingValidator
	authConfig       authservice.Config
	authorizationSvc authorizationservice.Service
}

func New(matchingSvc matchingservice.Service, authSvc authservice.Service, presenceSvc presenceservice.Service,
	validator matchingvalidator.MatchingValidator, authConfig authservice.Config) Handler {
	return Handler{
		matchingSvc: matchingSvc,
		authSvc:     authSvc,
		presenceSvc: presenceSvc,
		validator:   validator,
		authConfig:  authConfig,
	}
}
