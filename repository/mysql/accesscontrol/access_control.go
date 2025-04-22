package accesscontrol

import (
	"strings"

	"github.com/samber/lo"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

// TODO: two GetPermissionTitles and GetPermissionIDs are so similiar in syntax, check if you can perform better there

type PriviledgedType string

const (
	UserPriviledgedType  = "user"
	AdminPriviledgedType = "admin"
)

type SourcePemission struct {
	ID          uint8
	Title       string
	Description string
	CreatedAt   []uint8
}

type AccessControl struct {
	ID                 uint
	PrivilegedType     PriviledgedType
	PrivilegedID       uint
	SourcePermissionID uint
	CreatedAt          []uint8
}

func (s Storage) GetAllPermissionTitles(userID uint, userRole entity.Role) ([]string, error) {
	const op = "accesscontrol.GetAllPermissions"

	// With regard to performance and cache optimization, I decided to separate the db query
	// to persist the role query in db without binding it to the user query whcih makes it specific.
	userQuery := `SELECT source_permission_id FROM access_controls WHERE privileged_type = 'user' && privileged_id = ?`
	userPermissionIDs, uErr := s.getPermissionIDs(userQuery, userID)
	if uErr != nil {
		return nil, richerr.New(op).WithError(uErr)
	}

	// TODO: any explicit way to define longer cache persistency duration?
	roleQuery := `SELECT source_permission_id FROM access_controls WHERE privileged_type = 'role' && privileged_id = ?`
	rolePermissionIDs, rErr := s.getPermissionIDs(roleQuery, userRole)
	if rErr != nil {
		return nil, richerr.New(op).WithError(uErr)
	}

	allPermissionIDs := userPermissionIDs
	allPermissionIDs = append(allPermissionIDs, rolePermissionIDs...)

	if len(allPermissionIDs) < 1 {
		return []string{}, richerr.New(op).WithCode(richerr.ErrUnauthorized).WithMessage(errmsg.ErrMsgUnauthorized)
	}

	permissionTitles, pErr := s.GetPermissionTitles(lo.Uniq(allPermissionIDs)...)
	if pErr != nil {
		return nil, richerr.New(op).WithError(uErr)
	}

	return permissionTitles, nil
}

func (s Storage) getPermissionIDs(query string, previleged_id any) ([]int, error) {
	const op = "accesscontrol.getPermissionIDs"
	var fetchedAC AccessControl

	errFunc := func(wrappedError error) error {
		return richerr.New(op).WithError(wrappedError).
			WithCode(richerr.ErrUnexpected).WithMessage(errmsg.ErrMsgUnexpected)
	}

	rows, qErr := s.db.Query(query, previleged_id)
	if qErr != nil {
		return []int{}, errFunc(qErr)
	}
	defer rows.Close()

	permissionIDs := []int{}

	for rows.Next() {
		if sErr := rows.Scan(&fetchedAC.SourcePermissionID); sErr != nil {
			return []int{}, errFunc(sErr)
		}

		permissionIDs = append(permissionIDs, int(fetchedAC.SourcePermissionID))
	}

	if rows.Err() != nil {
		return []int{}, errFunc(rows.Err())
	}

	return permissionIDs, nil
}

func (s Storage) GetPermissionTitles(permissionIDs ...int) ([]string, error) {
	const op = "accesscontrol.GetPermissionTitles"

	errFunc := func(wrappedError error) error {
		return richerr.New(op).WithError(wrappedError).
			WithCode(richerr.ErrUnexpected).WithMessage(errmsg.ErrMsgUnexpected)
	}

	query := `SELECT title FROM source_permissions WHERE id IN(?` +
		strings.Repeat(", ?", len(permissionIDs)-1) + ")"

	args := make([]any, len(permissionIDs))
	for i, id := range permissionIDs {
		args[i] = id
	}

	rows, qErr := s.db.Query(query, args...)
	if qErr != nil {
		return []string{}, errFunc(qErr)
	}
	defer rows.Close()

	permissionTitles := []string{}
	for rows.Next() {
		var fetchedPermission SourcePemission
		if sErr := rows.Scan(&fetchedPermission.Title); sErr != nil {
			return []string{}, errFunc(sErr)
		}

		permissionTitles = append(permissionTitles, fetchedPermission.Title)
	}

	if rows.Err() != nil {
		return []string{}, errFunc(rows.Err())
	}

	return permissionTitles, nil
}
