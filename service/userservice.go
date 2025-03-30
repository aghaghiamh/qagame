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
}

type UserService struct {
	repo UserRepo
}

func New(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
	Password    string
}

type RegisterResponse struct {
	User entity.User
}

func (s *UserService) Register(req RegisterRequest) (RegisterResponse, error) {
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

	return RegisterResponse{User: user}, nil
}

// Validate phone number using (+98) 09xxxxxxxxx pattern which x is a digit
func isValidPhoneNumber(phoneNumber string) (bool, error) {
	r, err := regexp.Compile(`^(\(?\+98\)?)?[-\s]?(09)(\d{9})$`)
	if err != nil {
		return false, fmt.Errorf(`regexp pattern format is incorrect`)
	}

	if !r.MatchString(phoneNumber) {
		return false, fmt.Errorf(`phone number must start with "09", and all characters must be digits`)
	}

	return true, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &utils.RichErr{
			Kind:    utils.ServerError,
			Message: fmt.Sprintf("Couldn't hash the password: %s", password),
		}
	}
	return string(hashedPassword), nil
}
