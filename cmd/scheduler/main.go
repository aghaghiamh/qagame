package main

import (
	"github.com/aghaghiamh/gocast/QAGame/scheduler"

	"os"
	"os/signal"
	"time"
)

func main() {
	// run the cronjob scheduler
	schDoneCH := make(chan bool)
	sch := scheduler.New(schDoneCH)
	go func() {
		sch.Start()
	}()

	// Graceful Termination - wait until there is a os.signal on the quit channel then revoke all other children.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	schDoneCH <- true

	time.Sleep(3 * time.Second)
}
