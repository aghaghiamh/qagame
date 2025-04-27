package scheduler

import (
	"context"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/service/matchingservice"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	done chan bool
	matchingSVC matchingservice.Service
}

func New(done chan bool, matchingSVC matchingservice.Service) Scheduler {
	return Scheduler{
		done: done,
		matchingSVC: matchingSVC,
	}
}

func (s Scheduler) Start() {
	// a long-running func which has been registered with many other funcs with different periodic functionalieties
	// which have to take care of their run
	fmt.Println("Scheduler has been started!!!")

	c := cron.New()
	if _, err := c.AddFunc("* * * * *", s.ScheduleMatchPlayersInWaitingList); err != nil {
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