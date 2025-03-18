package entity

import "time"

type User struct {
	ID               uint
	Name             string
	PhoneNumber      string
	BirthDate        time.Time
	RegistrationDate time.Time
}
