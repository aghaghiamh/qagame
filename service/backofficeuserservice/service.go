package backofficeuserservice

import "github.com/aghaghiamh/gocast/QAGame/entity"

type Service struct {
}

func New() Service {
	return Service{}
}

// Warning: Mocked Behaviour
func (s Service) ListAllUsers() ([]entity.User, error) {
	return []entity.User{
		{
			ID:          2,
			Name:        "Aghaghia",
			PhoneNumber: "09027307801",
		},
	}, nil
}
