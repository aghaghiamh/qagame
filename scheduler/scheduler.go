package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
	done chan bool
}

func New(done chan bool) Scheduler {
	return Scheduler{
		done: done,
	}
}

func (s Scheduler) Start() {
	// a long-running func which has been registered with many other funcs with different periodic functionalieties
	// which have to take care of their run
	fmt.Println("Scheduler has been started!!!")
	for {
		select {
		case <-s.done:
			fmt.Println("Scheduler has been shutdowned gracefully!!!")
			return
		default:
			fmt.Println(time.Now())
			time.Sleep(1 * time.Second)
		}
	}
}