package entity

type PermissionTitle string

const (
	UserListPermission = PermissionTitle("user-list")
)

type SourcePemission struct {
	ID          uint8
	Title       string
	Description string
}
