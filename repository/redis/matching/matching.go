package matching

import (
	"context"
	"strconv"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/aghaghiamh/gocast/QAGame/pkg/timestamp"
	"github.com/redis/go-redis/v9"
)

func (s *Storage) AddToWaitingList(ctx context.Context, key string, userID uint) error {
	const op = richerr.Operation("redis.matching.AddToWaitingList")

	// As if the user already requests for a game in the specific category, we just update it's timestamp to
	// Now, we might get the 0 for number of Operation have happened, then it shouldn't be taken into account.
	_, aErr := s.adapter.Driver().ZAdd(ctx, key, redis.Z{
		Score:  float64(timestamp.Now()),
		Member: userID,
	}).Result()

	if aErr != nil {
		return richerr.New(op).WithCode(richerr.ErrUnexpected).WithMessage(errmsg.ErrMsgUnexpected)
	}

	return nil
}

func (s *Storage) GetFromWaitingList(ctx context.Context, key string, maxNumOfUsers int) ([]entity.WaitingMember, error) {
	const op = richerr.Operation("redis.matching.GetFromWaitingList")

	// We can also get the timestamp with the ZRangeWithScores but it seems redundant and not used at the moment.
	rangeRes, rErr := s.adapter.Driver().ZRange(ctx, key, 0, int64(maxNumOfUsers)-1).Result()
	if rErr != nil {
		return []entity.WaitingMember{}, richerr.New(op).WithError(rErr).
			WithCode(richerr.ErrUnexpected).WithMessage(errmsg.ErrMsgUnexpected)
	}

	WatingMems := []entity.WaitingMember{}
	for _, userIDstr := range rangeRes {
		userID, _ := strconv.Atoi(userIDstr)
		WatingMems = append(WatingMems, entity.WaitingMember{
			UserID: uint(userID),
		})
	}

	return WatingMems, nil
}

func (s *Storage) RemoveFromWaitingList(ctx context.Context, key string, userIDs []uint) error {
	const op = richerr.Operation("redis.matching.RemoveFromWaitingList")

	// Convert uint slice to interface{} slice for Redis
	members := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		members[i] = strconv.FormatUint(uint64(id), 10)
	}

	if _, err := s.adapter.Driver().ZRem(ctx, key, members...).Result(); err != nil {
		return richerr.New(op).WithError(err).WithCode(richerr.ErrUnexpected)
	}

	return nil
}
