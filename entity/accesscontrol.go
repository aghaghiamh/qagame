package entity

type PriviledgedType string

const (
	UserPriviledgedType  = "user"
	AdminPriviledgedType = "admin"
)

type AccessControl struct {
	ID                 uint
	PriviledgedType    PriviledgedType
	PriviledgedID      uint
	SourcePermissionID uint
}
