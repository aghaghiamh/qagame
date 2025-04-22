package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

func MapToEntityRole(role string) Role {
	switch role {
	case AdminPriviledgedType:
		return AdminRole
	default:
		return UserRole
	}
}

func MapToEnumRole(role Role) string {
	switch role {
	case AdminRole:
		return AdminPriviledgedType
	default:
		return UserPriviledgedType
	}
}
