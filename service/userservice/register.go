package userservice

import (
	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

func (s *Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	const op = "userservice.Register"

	// Determine whether it's a normal path of password validation, or an OTP one.
	var hashedPassword string
	if len(req.Password) == 0 {
		// TODO: OTP path of Registeration, Currently it would recognized as a abnormality in validation logic
	} else {
		var passErr error
		hashedPassword, passErr = hashPassword(req.Password)
		if passErr != nil {

			return dto.RegisterResponse{}, richerr.New(op).WithError(passErr)
		}
	}

	user := entity.User{
		PhoneNumber:    req.PhoneNumber,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}
	var regErr error
	user, regErr = s.repo.Register(user)
	if regErr != nil {
		return dto.RegisterResponse{}, richerr.New(op).WithError(regErr)
	}

	return dto.RegisterResponse{
		UserInfo: dto.UserInfo{
			UserID:      user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}
