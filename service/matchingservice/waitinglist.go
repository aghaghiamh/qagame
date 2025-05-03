package matchingservice

import (
	"context"
	"fmt"
	"sync"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
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
	matchedUsers := []entity.MatchedUsers{}
	var mutex sync.Mutex

	for _, category := range entity.AllCategories() {
		wg.Add(1)
		go func(cat string) {
			defer wg.Done()

			categoryMatchedPlayers, err := s.matchCategoryPalyers(ctx, cat)
			if err != nil {
				richerr.New(op).WithError(err)
				return
			}

			mutex.Lock()
			matchedUsers = append(matchedUsers, categoryMatchedPlayers...)
			mutex.Unlock()
		}(string(category))
	}

	wg.Wait()
	return dto.MatchPlayersResponse{
		MatchedUsers: matchedUsers,
	}, nil
}

func (s Service) matchCategoryPalyers(ctx context.Context, category string) ([]entity.MatchedUsers, error) {
	const op = richerr.Operation("matchingservice.matchCategoryPalyers")

	key := s.genWaitingListKey(string(category))
	wMems, wErr := s.repo.GetFromWaitingList(ctx, key, s.config.maxNumOfUsers)
	if wErr != nil {
		return []entity.MatchedUsers{}, richerr.New(op).WithError(wErr)
	}

	presenceReq := dto.PresenceGetUsersInfoRequest{
		UserIDs: lo.Map(wMems, func(wMem entity.WaitingMember, _ int) uint {
			return wMem.UserID
		}),
	}
	resp, err := s.presenceClient.GetUsersAvailabilityInfo(ctx, presenceReq)
	if err != nil {
		return []entity.MatchedUsers{}, richerr.New(op).WithError(err)
	}

	removeFromRedis := []uint{}
	onlineUsers := []entity.UserAvailabilityInfo{}
	for _, uInfo := range resp.UsersAvailabilityInfo {
		if !lo.Contains(presenceReq.UserIDs, uInfo.UserID) {
			removeFromRedis = append(removeFromRedis, uInfo.UserID)
			continue
		}
		onlineUsers = append(onlineUsers, uInfo)
	}

	// TODO: remove the offline players from redis
	// if err := s.repo.RemoveFromWaitingList(ctx, key, removeFromRedis); err != nil {
	// 	fmt.Println(err)
	// 	return [][]uint{}, richerr.New(op).WithError(err)
	// }

	matchedUserIDs := []entity.MatchedUsers{}
	for i := 0; i < len(onlineUsers)-1; i += 2 {
		matchedUserIDs = append(matchedUserIDs,
			entity.MatchedUsers{
				Category: entity.Category(category),
				UserIDs:  []uint{onlineUsers[i].UserID, onlineUsers[i+1].UserID},
			})

		// TODO: publish an event here for matched players, and do not create a matchedUserIDs to pass to the parent func.
		// TODO: remove or add to removeFromRedis list to be removed from the zset key.
	}

	return matchedUserIDs, nil
}

func (s Service) genWaitingListKey(category string) string {
	// prefix:category
	return fmt.Sprint(s.config.RedisWaitingPrefix, ":", category)
}
