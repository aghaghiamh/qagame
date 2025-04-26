package matchingservice

import (
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
