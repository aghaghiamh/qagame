package presenceservice

import (
	"context"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/samber/lo"
)

func (s Service) GetUsersAvailabilityInfo(ctx context.Context, req dto.PresenceGetUsersInfoRequest) (dto.PresenceGetUsersInfoResponse, error) {
	const op = richerr.Operation("presenceservice.GetUsersAvailabilityInfo")

	keys := lo.Map(req.UserIDs, func(u uint, _ int) string {
		return s.generateKey(u)
	})

	usersTimestamp, err := s.repo.GetUsersTimestamp(ctx, keys)
	if err != nil {
		return dto.PresenceGetUsersInfoResponse{}, richerr.New(op).WithError(err)
	}

	usersAvailabilityInfo := lo.Map(req.UserIDs, func(uID uint, i int) entity.UserAvailabilityInfo {
		return entity.UserAvailabilityInfo{
			UserID:       uID,
			LastOnlineAt: usersTimestamp[i],
		}
	})

	return dto.PresenceGetUsersInfoResponse{
		UsersAvailabilityInfo: usersAvailabilityInfo,
	}, nil
}
