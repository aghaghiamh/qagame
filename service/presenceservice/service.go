package presenceservice

import (
	"context"
	"fmt"
	"time"
)

type Config struct {
	ExpectedOnlineTime time.Duration `mapstructure:"expected_online_time"`
	Prefix             string        `mapstructure:"prefix"`
}

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, ttl time.Duration) error
	GetUsersTimestamp(ctx context.Context, keys []string) ([]int64, error)
}

type Service struct {
	config Config
	repo   Repository
}

func New(config Config, repo Repository) Service {
	return Service{
		config: config,
		repo:   repo,
	}
}

func (s Service) generateKey(userID uint) string {
	return fmt.Sprintf("%s:%d", s.config.Prefix, userID)
}
