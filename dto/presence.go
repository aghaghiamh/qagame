package dto

import (
	"github.com/aghaghiamh/gocast/QAGame/entity"
)

type PresenceUpsertRequest struct {
	UserID    uint
	Timestamp int64
}

type PresenceUpsertResponse struct {
}

type PresenceGetUsersInfoRequest struct {
	UserIDs []uint
}

type PresenceGetUsersInfoResponse struct {
	UsersAvailabilityInfo []entity.UserAvailabilityInfo
}
