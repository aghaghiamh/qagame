package matchingservice

import "time"

type Service struct {
	repo      Reopository
	svcConfig MatchingServiceConfig
}

type MatchingServiceConfig struct {
	WaitingTimeout     time.Duration `mapstructure:"waiting_timeout"`
	RedisWaitingPrefix string        `mapstructure:"waiting_prefix"`
}

type Reopository interface {
	AddToWaitingList(string, uint) error
}

func New(repo Reopository, svcConfig MatchingServiceConfig) Service {
	return Service{
		repo:      repo,
		svcConfig: svcConfig,
	}
}
