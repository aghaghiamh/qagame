package matching

import (
	"context"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) AddToWaitingList(key string, userID uint) error {
	const op = "redis.matching.AddToWaitingList"

	// As if the user already requests for a game in the specific category, we just update it's timestamp to
	// Now, we might get the 0 for number of Operation have happened, then it shouldn't be taken into account.
	_, aErr := s.adapter.Driver().ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: userID,
	}).Result()

	if aErr != nil {
		return richerr.New(op).WithCode(richerr.ErrUnexpected).WithMessage(errmsg.ErrMsgUnexpected)
	}

	return nil
}

// func (s *Storage) GetFromWaitingList(key string) ([]uint, error) {
//  // Not is being used at the moment - dead code.
// 	const op = "redis.matching.AddToWaitingList"

// 	userIDsStr, rErr := s.adapter.Driver().ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
// 		Offset: 0,
// 		Count: time.Now().Unix(),
// 	}).Result()

// 	if rErr != nil {
// 		return []uint{}, richerr.New(op).WithCode(richerr.ErrUnexpected).WithMessage(errmsg.ErrMsgUnexpected)
// 	}

// 	userIDs := lo.Map(userIDsStr, func(str string, _ int) uint {
// 		num, _ := strconv.Atoi(str)
// 		return uint(num)
// 	})

// 	return userIDs, nil
// }
