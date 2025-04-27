package matchingservice

import (
	"context"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

func (s Service) AddToWaitingList(req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"
	// prefix:category {prefix is stored in the matchingSvcConfig and category is in the req payload}
	waitingKey := fmt.Sprint(s.svcConfig.RedisWaitingPrefix, ":", string(req.Category))

	// if the player is already in the category list, update its timestamp
	repoErr := s.repo.AddToWaitingList(waitingKey, req.UserID)
	if repoErr != nil {
		return dto.AddToWaitingListResponse{}, richerr.New(op).WithError(repoErr)
	}

	return dto.AddToWaitingListResponse{WaitingListTimeout: s.svcConfig.WaitingTimeout}, nil
}

func (s Service) MatchPlayers(ctx context.Context, req dto.MatchPlayersRequest) (dto.MatchPlayersResponse, error) {
	const op = "matchingservice.MatchPlayers"
	fmt.Println(op)
	// get the users in the provided category in req.Category from the redis
	// based on some logic (probably random at the moment) match each 2 players exist in the redis category key
	// if all players have been matched together successfuly, remove all from the redis in bulk using ZREMby..., otherway
	// (odd number or a fail in match) remove one by one from redis
	// return the pair of matched players in a list of tuples
	return dto.MatchPlayersResponse{}, nil
}
