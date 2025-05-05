package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aghaghiamh/gocast/QAGame/adapter/presenceclient"
	"github.com/aghaghiamh/gocast/QAGame/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := fmt.Sprintf(":%d", 8089)
	conn, cErr := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if cErr != nil {
		log.Fatalf("couldn't create the presence clinet on %s address.", address)
	}
	defer conn.Close()

	client := presenceclient.New(conn)
	resp, gErr := client.GetUsersAvailabilityInfo(context.Background(), dto.PresenceGetUsersInfoRequest{
		UserIDs: []uint{2},
	})
	if gErr != nil {
		log.Printf(gErr.Error())
		return
	}
	for _, uInfo := range resp.UsersAvailabilityInfo {
		fmt.Println("userID: ", uInfo.UserID, " user Last Online Timestamp: ", uInfo.LastOnlineAt)
	}
}
