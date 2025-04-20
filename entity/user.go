package entity

import "time"

type User struct {
	ID               uint
	Name             string
	PhoneNumber      string
	HashedPassword   string
	BirthDate        time.Time
	RegistrationDate time.Time
	Role             Role // for authorization
}
