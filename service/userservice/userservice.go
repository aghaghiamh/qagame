package userservice

import (
	"fmt"

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
	CreateAccessToken(userID uint, userRole entity.Role) (string, error)
	CreateRefreshToken(userID uint, userRole entity.Role) (string, error)
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
