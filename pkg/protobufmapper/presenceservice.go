package protobufmapper

import (
	"github.com/aghaghiamh/gocast/QAGame/contract/golang/presence"
	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
)

func MapProtoUserAvailabilityInfoResponseToDto(usersInfo presence.GetUsersAvailabilityInfoResponse) *dto.PresenceGetUsersInfoResponse {
	dtoUsersInfo := &dto.PresenceGetUsersInfoResponse{}
	for _, uInfo := range usersInfo.UsersInfo {
		dtoUsersInfo.UsersAvailabilityInfo = append(dtoUsersInfo.UsersAvailabilityInfo,
			entity.UserAvailabilityInfo{
				UserID:       uint(uInfo.UserId),
				LastOnlineAt: uInfo.LastOnlineAt,
			})
	}
	return dtoUsersInfo
}

func MapDtoUserAvailabilityInfoResponseToProto(dto dto.PresenceGetUsersInfoResponse) *presence.GetUsersAvailabilityInfoResponse {
	protoUsersInfo := &presence.GetUsersAvailabilityInfoResponse{}
	for _, uInfo := range dto.UsersAvailabilityInfo {
		protoUsersInfo.UsersInfo = append(protoUsersInfo.UsersInfo, &presence.UserAvailabilityInfo{
			UserId:       uint64(uInfo.UserID),
			LastOnlineAt: uInfo.LastOnlineAt,
		})
	}
	return protoUsersInfo
}
