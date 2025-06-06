package userservice

import (
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"

	var loginResponse dto.LoginResponse

	if len(req.Password) == 0 {
		// TODO: OTP Login
	} else { // Normal Authentication Path
		user, dbErr := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
		if dbErr != nil {
			return dto.LoginResponse{}, richerr.New(op).
				WithError(dbErr).
				WithMessage(fmt.Sprintf(
					"User with %s phone number does not exist. Please register first", req.PhoneNumber))
		}

		if len(user.HashedPassword) == 0 {
			// TODO: OTP Login
		}

		// verify password
		verErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
		if verErr != nil {
			return dto.LoginResponse{}, richerr.New(op).
				WithError(verErr).
				WithCode(richerr.ErrInvalidInput).
				// TODO: Should I hide the actual errror in this layer?
				WithMessage("Phone number or password is incorrect.")
		}

		accessToken, aErr := s.auth.CreateAccessToken(user.ID, user.Role)
		if aErr != nil {
			return dto.LoginResponse{}, richerr.New(op).WithError(aErr)
		}

		refreshToken, rErr := s.auth.CreateRefreshToken(user.ID, user.Role)
		if rErr != nil {
			return dto.LoginResponse{}, richerr.New(op).WithError(rErr)
		}

		loginResponse = dto.LoginResponse{
			UserInfo: dto.UserInfo{
				UserID:      user.ID,
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
			},
			AuthTokens: dto.AuthTokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}
	}
	return loginResponse, nil
}
