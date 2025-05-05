package matchingservice

import (
	"context"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/contract/broker"
	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
)

type PresenceClient interface {
	GetUsersAvailabilityInfo(ctx context.Context, req dto.PresenceGetUsersInfoRequest) (dto.PresenceGetUsersInfoResponse, error)
}

type Service struct {
	repo           Reopository
	config         Config
	broker         broker.Broker
	presenceClient PresenceClient
}

type Config struct {
	maxNumOfUsers      int           `mapstructure:"max_num_of_users_to_be_fetched_in_each_iter"`
	WaitingTimeout     time.Duration `mapstructure:"waiting_timeout"`
	RedisWaitingPrefix string        `mapstructure:"waiting_prefix"`
}

// For Async/some-other use-cases ctx should be initialized inside the funciton.
type Reopository interface {
	AddToWaitingList(ctx context.Context, key string, userID uint) error
	GetFromWaitingList(ctx context.Context, key string, maxNumOfUsers int) ([]entity.WaitingMember, error)
	RemoveFromWaitingList(key string, userIDs []uint)
}

func New(repo Reopository, config Config, broker broker.Broker, presenceClient PresenceClient) Service {
	return Service{
		repo:           repo,
		config:         config,
		broker:         broker,
		presenceClient: presenceClient,
	}
}
