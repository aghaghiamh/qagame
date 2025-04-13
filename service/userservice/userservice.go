package userservice

import (
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Register(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	GetUserByID(user_id uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(userID uint) (string, error)
	CreateRefreshToken(userID uint) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo UserRepo
}

func New(repo UserRepo, auth AuthGenerator) Service {
	return Service{
		auth: auth,
		repo: repo,
	}
}

func (s *Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"

	// if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {
	// 	return dto.LoginResponse{}, err
	// }

	var loginResponse dto.LoginResponse

	if len(req.Password) == 0 {
		// TODO: OTP Login
	} else if len(req.Password) < 8 {

		return dto.LoginResponse{}, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage("Password should contains at least 8 characters")
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

		accessToken, aErr := s.auth.CreateAccessToken(user.ID)
		if aErr != nil {
			return dto.LoginResponse{}, richerr.New(op).WithError(aErr)
		}

		refreshToken, rErr := s.auth.CreateRefreshToken(user.ID)
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

func (s *Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	const op = "userservice.Register"

	// password validation
	var hashedPassword string
	if len(req.Password) == 0 {
		// TODO: OTP path of Registeration
	} else if len(req.Password) < 8 {

		return dto.RegisterResponse{}, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage("user's password must contains at least 8 character")
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

func hashPassword(password string) (string, error) {
	const op = "userservice.hashPassword"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return "", richerr.New(op).
			WithCode(richerr.ErrServer).
			WithMessage(fmt.Sprintf("Couldn't hash the password: %s", password))
	}

	return string(hashedPassword), nil
}

func (s *Service) GetProfile(req dto.UserProfileRequest) (dto.UserProfileResponse, error) {
	const op = "userservice.GetUserProfile"
	user, gErr := s.repo.GetUserByID(req.UserID)
	if gErr != nil {

		return dto.UserProfileResponse{}, richerr.New(op).
			WithError(gErr).
			WithMessage(fmt.Sprintf("User with %d id does not exist. Please register first.", req.UserID))
	}

	return dto.UserProfileResponse{
		UserInfo: dto.UserInfo{
			UserID:      user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}
