package uservalidator

type UserRepo interface {
	IsAlreadyExist(phoneNumber string) (bool, error)
}

type UserValidator struct {
	repo UserRepo
}

func New(repo UserRepo) UserValidator {
	return UserValidator{
		repo: repo,
	}
}
