package uservalidator

import "github.com/aghaghiamh/gocast/QAGame/entity"

const (
	PhoneNumberRegex = `^(\(?\+98\)?)?[-\s]?(09)(\d{9})$`
)

type UserRepo interface {
	IsAlreadyExist(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type UserValidator struct {
	repo UserRepo
}

func New(repo UserRepo) UserValidator {
	return UserValidator{
		repo: repo,
	}
}
