package gutils

import (
	"time"
)

// Schedule schedules a task that is running in the background and can be kill by the chan returned by the function
func Schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
