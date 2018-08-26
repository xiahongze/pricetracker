package gutils

import (
	"log"
	"os"
	"strconv"
	"time"
)

// SchdFreq is the default scheduled frequency
var SchdFreq = time.Minute * 1

func init() {
	if v, ok := os.LookupEnv("SCHEDULE_FREQ"); ok {
		if vi, err := strconv.Atoi(v); err == nil {
			SchdFreq = time.Minute * time.Duration(vi)
		} else {
			log.Println("WARN: GUTILS: SCHEDULE_FREQ is not int but ", v)
		}
	}
}

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
