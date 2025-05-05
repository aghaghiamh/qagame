package presenceclient

import (
	"context"

	"github.com/aghaghiamh/gocast/QAGame/contract/goproto/presence"
	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/protobufmapper"
	"github.com/aghaghiamh/gocast/QAGame/pkg/typemapper"
	"google.golang.org/grpc"
)

type Client struct {
	client presence.PresenceServiceClient
}

func New(conn *grpc.ClientConn) Client {
	return Client{
		client: presence.NewPresenceServiceClient(conn),
	}
}

func (c Client) GetUsersAvailabilityInfo(ctx context.Context, req dto.PresenceGetUsersInfoRequest) (
	dto.PresenceGetUsersInfoResponse, error) {

	userIDs := typemapper.ArrayMapper(req.UserIDs, func(uID uint) uint64 {
		return uint64(uID)
	})

	resp, err := c.client.GetUsersAvailabilityInfo(ctx, &presence.GetUsersAvailabilityInfoRequest{
		UserIds: userIDs})

	if err != nil {
		return dto.PresenceGetUsersInfoResponse{}, err
	}

	return *protobufmapper.MapProtoUserAvailabilityInfoResponseToDto(*resp), nil
}
