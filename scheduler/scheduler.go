package scheduler

import (
	"context"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/service/matchingservice"
	"github.com/robfig/cron/v3"
)

type Config struct {
	MatchPlayersCronjobIntervalsInMins string `mapstructure:"match_players_cronjob_intervals_in_mins"`
}

type Scheduler struct {
	config      Config
	done        chan bool
	matchingSVC matchingservice.Service
}

func New(config Config, done chan bool, matchingSVC matchingservice.Service) Scheduler {
	return Scheduler{
		config:      config,
		done:        done,
		matchingSVC: matchingSVC,
	}
}

func (s Scheduler) Start() {
	// a long-running func which has been registered with many other funcs with different periodic functionalieties
	// which have to take care of their run
	fmt.Println("Scheduler has been started!!!")

	c := cron.New()
	matchPlayersIntervals := fmt.Sprintf("%s * * * *", s.config.MatchPlayersCronjobIntervalsInMins)
	if _, err := c.AddFunc(matchPlayersIntervals, s.ScheduleMatchPlayersInWaitingList); err != nil {
		fmt.Println("Schedule Err: ", err)
	}
	c.Start()

	select {
	case <-s.done:
		fmt.Println("Scheduler has been shutdowned gracefully!!!")
		return
	}
}

func (s Scheduler) ScheduleMatchPlayersInWaitingList() {
	fmt.Println("ScheduleMatchPlayersInWaitingList")
	s.matchingSVC.MatchPlayers(context.Background(), dto.MatchPlayersRequest{})
}
