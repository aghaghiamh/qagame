package userservice

import (
	"fmt"
	"regexp"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	utils "github.com/aghaghiamh/gocast/QAGame/pkg"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	IsAlreadyExist(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	GetUser(phoneNumber string) (entity.User, error)
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

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserInfo struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterResponse struct {
	UserInfo
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRequest struct {
	UserID       uint
	AccessToken  string
	RefreshToken string
}

type AuthResponse struct {
	AccessToken string
}

func (s *Service) Login(req LoginRequest) (LoginResponse, error) {
	if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {
		return LoginResponse{}, err
	}

	var loginResponse LoginResponse

	if len(req.Password) == 0 {
		// TODO: OTP Login
	} else if len(req.Password) < 8 {
		return LoginResponse{}, &utils.RichErr{
			Code:    utils.GeneralBusinessLogicError,
			Message: "Password should contains at least 8 characters",
		}
	} else { // Normal Authentication Path
		user, err := s.repo.GetUser(req.PhoneNumber)
		if err != nil {
			return LoginResponse{}, err
		}

		if len(user.HashedPassword) == 0 {
			// TODO: OTP Login
		}

		// verify password
		verErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
		if verErr != nil {
			return LoginResponse{}, utils.RichErr{
				Code:    utils.GeneralBusinessLogicError,
				Message: "Phone number or password is incorrect.",
			}
		}

		accessToken, aErr := s.auth.CreateAccessToken(user.ID)
		if aErr != nil {
			return LoginResponse{}, aErr
		}

		refreshToken, rErr := s.auth.CreateRefreshToken(user.ID)
		if rErr != nil {
			return LoginResponse{}, rErr
		}

		loginResponse = LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
	}
	return loginResponse, nil
}

func (s *Service) Register(req RegisterRequest) (RegisterResponse, error) {
	if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {
		return RegisterResponse{}, err
	}

	alreadyExist, err := s.repo.IsAlreadyExist(req.PhoneNumber)
	if alreadyExist {
		return RegisterResponse{}, fmt.Errorf("duplicated user with %s phone number", req.PhoneNumber)
	}
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("database error: %w", err)
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("user's name must contains at least 3 character")
	}

	// password validation
	var hashedPassword string
	if len(req.Password) == 0 {
		// TODO: OTP path of Registeration
	} else if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("user's password must contains at least 8 character")
	} else {
		var passErr error
		hashedPassword, passErr = hashPassword(req.Password)
		if passErr != nil {
			return RegisterResponse{}, passErr
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
		return RegisterResponse{}, regErr
	}

	return RegisterResponse{
		UserInfo{
			UserID:      user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}

// Validate phone number using (+98) 09xxxxxxxxx pattern which x is a digit
func isValidPhoneNumber(phoneNumber string) (bool, error) {
	r, err := regexp.Compile(`^(\(?\+98\)?)?[-\s]?(09)(\d{9})$`)
	if err != nil {
		return false, &utils.RichErr{
			Code:    utils.GeneralBusinessLogicError,
			Message: "regexp pattern format is incorrect",
		}
	}

	if !r.MatchString(phoneNumber) {
		return false, &utils.RichErr{
			Code:    utils.GeneralBusinessLogicError,
			Message: `phone number must start with "09", and all characters must be digits`,
		}
	}

	return true, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &utils.RichErr{
			Code:    utils.GeneralServerErr,
			Message: fmt.Sprintf("Couldn't hash the password: %s", password),
		}
	}
	return string(hashedPassword), nil
}
