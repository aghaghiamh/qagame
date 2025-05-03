package presence

import (
	"context"
	"strconv"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"github.com/samber/lo"
)

func (s Storage) Upsert(ctx context.Context, key string, timestamp int64, ttl time.Duration) error {
	const op = richerr.Operation("redis.presence.Upsert")

	_, sErr := s.adapter.Driver().Set(ctx, key, timestamp, ttl).Result()
	if sErr != nil {
		return richerr.New(op).WithCode(richerr.ErrUnexpected).WithError(sErr)
	}

	return nil
}

func (s Storage) GetUsersTimestamp(ctx context.Context, keys []string) ([]int64, error) {
	const op = richerr.Operation("redis.presence.GetUsersTimestamp")

	userTsStrs, err := s.adapter.Driver().MGet(ctx, keys...).Result()
	if err != nil {
		return []int64{}, richerr.New(op).WithError(err).WithCode(richerr.ErrUnexpected).
			WithMessage(errmsg.ErrMsgUnexpected)
	}

	listUserTs := lo.Map(userTsStrs, func(ts interface{}, _ int) int64 {
		userTs, _ := strconv.Atoi(ts.(string))
		return int64(userTs)
	})

	return listUserTs, nil
}
