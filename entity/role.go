package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)
