package authorizationservice

import (
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/samber/lo"
)

type Repository interface {
	GetAllPermissionTitles(userID uint, userRole entity.Role) ([]string, error)
}

type Service struct {
	acRepo Repository
}

func New(access_control_repo Repository) Service {
	return Service{
		acRepo: access_control_repo,
	}
}

func (s Service) CheckPermissions(userID uint, userRole entity.Role, requiredPermissionTitles ...entity.PermissionTitle) (bool, error) {
	const op = "authenticationservice.CheckPermissions"

	userPermissionTitles, repoErr := s.acRepo.GetAllPermissionTitles(userID, userRole)
	if repoErr != nil {
		return false, richerr.New(op)
	}

	// Convert user permissions to a map for O(1) lookups
	userPermissionMap := lo.SliceToMap(userPermissionTitles, func(perm string) (entity.PermissionTitle, bool) {
		return entity.PermissionTitle(perm), true
	})

	// Check if all required permissions exist in the user's permissions
	hasAllPermissions := lo.EveryBy(requiredPermissionTitles, func(requiredPerm entity.PermissionTitle) bool {

		_, exists := userPermissionMap[requiredPerm]
		return exists
	})

	return hasAllPermissions, nil
}
