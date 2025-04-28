package presenceservice

import (
	"context"
	"time"
)

type Config struct {
	TTL    time.Duration `mapstructure:"ttl"`
	Prefix string        `mapstructure:"prefix"`
}

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, ttl time.Duration) error
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
