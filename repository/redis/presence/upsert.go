package presence

import (
	"context"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

func (s Storage) Upsert(ctx context.Context, key string, timestamp int64, ttl time.Duration) error {
	const op = richerr.Operation("redis.presence.Upsert")

	_, sErr := s.adapter.Driver().Set(ctx, key, timestamp, ttl).Result()
	if sErr != nil {
		return richerr.New(op).WithCode(richerr.ErrUnexpected).WithError(sErr)
	}

	return nil
}
