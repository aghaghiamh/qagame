package matchingservice

import "time"

type Service struct {
	repo   Reopository
	config Config
}

type Config struct {
	WaitingTimeout     time.Duration `mapstructure:"waiting_timeout"`
	RedisWaitingPrefix string        `mapstructure:"waiting_prefix"`
}

type Reopository interface {
	AddToWaitingList(string, uint) error
}

func New(repo Reopository, config Config) Service {
	return Service{
		repo:   repo,
		config: config,
	}
}
