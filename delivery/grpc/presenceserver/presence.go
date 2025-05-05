package presenceserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aghaghiamh/gocast/QAGame/contract/goproto/presence"
	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/protobufmapper"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/aghaghiamh/gocast/QAGame/pkg/typemapper"
	"github.com/aghaghiamh/gocast/QAGame/service/presenceservice"
	"google.golang.org/grpc"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	Service *presenceservice.Service
}

func New(service *presenceservice.Service) Server {
	return Server{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		Service:                            service,
	}
}

func (s Server) GetUsersAvailabilityInfo(ctx context.Context, req *presence.GetUsersAvailabilityInfoRequest) (
	*presence.GetUsersAvailabilityInfoResponse, error) {
	const op = richerr.Operation("presenceservice.GetUsersAvailabilityInfo")

	log.Printf("%s - user IDs to be fetched: %v\n", op, req.UserIds)

	userIDs := typemapper.ArrayMapper(req.UserIds, func(u uint64) uint { return uint(u) })
	resp, iErr := s.Service.GetUsersAvailabilityInfo(ctx, dto.PresenceGetUsersInfoRequest{
		UserIDs: userIDs,
	})
	if iErr != nil {
		fmt.Println(iErr)
		// TODO: write a grpc mapper for richerr(s) like httpmapper package in pkg
		return nil, iErr
	}

	return protobufmapper.MapDtoUserAvailabilityInfoResponseToProto(resp), nil
}

func (s Server) Serve() {
	address := fmt.Sprintf(":%d", 8089)
	listener, lErr := net.Listen("tcp", address)
	if lErr != nil {
		log.Fatalf("can't listen on %s address", address)
	}

	grpcServer := grpc.NewServer()
	presence.RegisterPresenceServiceServer(grpcServer, &s)

	log.Println("presence grpc server starting on", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("couldn't serve presence grpc server")
	}
}
