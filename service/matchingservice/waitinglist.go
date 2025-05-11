package matchingservice

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/eventencoder"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/samber/lo"
)

func (s Service) AddToWaitingList(ctx context.Context, req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"

	waitingKey := s.genWaitingListKey(string(req.Category))

	// if the player is already in the category list, update its timestamp
	repoErr := s.repo.AddToWaitingList(ctx, waitingKey, req.UserID)
	if repoErr != nil {
		return dto.AddToWaitingListResponse{}, richerr.New(op).WithError(repoErr)
	}

	return dto.AddToWaitingListResponse{WaitingListTimeout: s.config.WaitingTimeout}, nil
}

func (s Service) MatchPlayers(ctx context.Context, _ dto.MatchPlayersRequest) (dto.MatchPlayersResponse, error) {
	const op = richerr.Operation("matchingservice.MatchPlayers")

	wg := &sync.WaitGroup{}

	for _, category := range entity.AllCategories() {
		wg.Add(1)
		go func(cat string) {
			defer wg.Done()

			err := s.matchCategoryPalyers(ctx, cat)
			if err != nil {
				richerr.New(op).WithError(err)
				return
			}

		}(string(category))
	}

	wg.Wait()
	return dto.MatchPlayersResponse{}, nil
}

func (s Service) matchCategoryPalyers(ctx context.Context, category string) error {
	const op = richerr.Operation("matchingservice.matchCategoryPalyers")

	key := s.genWaitingListKey(string(category))
	wMems, wErr := s.repo.GetFromWaitingList(ctx, key, s.config.maxNumOfUsers)
	if wErr != nil {
		return richerr.New(op).WithError(wErr)
	}

	if len(wMems) <= 0 {
		log.Printf("No waited user for %s key", key)
		return nil
	}

	presenceReq := dto.PresenceGetUsersInfoRequest{
		UserIDs: lo.Map(wMems, func(wMem entity.WaitingMember, _ int) uint {
			return wMem.UserID
		}),
	}
	resp, err := s.presenceClient.GetUsersAvailabilityInfo(ctx, presenceReq)
	if err != nil {
		return richerr.New(op).WithError(err)
	}

	uIDsToBeRemvoed := []uint{}
	onlineUsers := []entity.UserAvailabilityInfo{}
	for _, uInfo := range resp.UsersAvailabilityInfo {
		if !lo.Contains(presenceReq.UserIDs, uInfo.UserID) {
			uIDsToBeRemvoed = append(uIDsToBeRemvoed, uInfo.UserID)
			continue
		}
		onlineUsers = append(onlineUsers, uInfo)
	}

	go s.repo.RemoveFromWaitingList(key, uIDsToBeRemvoed)

	matchedUserIDsTobeRemoved := []uint{}
	for i := 0; i < len(onlineUsers)-1; i += 2 {
		matchedPlayersIDs := []uint{onlineUsers[i].UserID, onlineUsers[i+1].UserID}
		matchedPlayers := entity.MatchedPlayers{
			Category: entity.Category(category),
			UserIDs:  matchedPlayersIDs}
		matchedUserIDsTobeRemoved = append(matchedUserIDsTobeRemoved, matchedPlayersIDs...)

		payload, err := eventencoder.MatchedPlayerUsersEncoder(matchedPlayers)
		if err != nil {
			return richerr.New(op).WithError(err)
		}
		s.broker.Publish(entity.MatchingMatchedUsersEvent, payload)
	}
	go s.repo.RemoveFromWaitingList(key, matchedUserIDsTobeRemoved)

	return nil
}

func (s Service) genWaitingListKey(category string) string {
	// prefix:category
	return fmt.Sprint(s.config.RedisWaitingPrefix, ":", category)
}
