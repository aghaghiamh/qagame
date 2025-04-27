package dto

import (
	"time"

	"github.com/aghaghiamh/gocast/QAGame/entity"
)

type AddToWaitingListRequest struct {
	UserID   uint
	Category entity.Category `json:"category"`
}

type AddToWaitingListResponse struct {
	WaitingListTimeout time.Duration `json:"waiting_list_timeout_in_nanoseconds"`
}

type MatchPlayersRequest struct {
	Category entity.Category `json:"category"`
}

type MatchPlayersResponse struct {
}
