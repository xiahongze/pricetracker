package gutils

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/pushover"
	"github.com/xiahongze/pricetracker/trackers"
)

var priceRegex, _ = regexp.Compile("\\d+\\.?\\d{0,}")

func processEntity(ent *models.Entity, pushClient *pushover.Client) (err error) {
	// save the entity before returning
	defer func() {
		if err != nil {
			log.Printf("ERROR: %v", err)
			key, _ := ent.K.MarshalJSON()
			log.Printf("INFO: URL: %s\tXPATH: %s\tKey: %s", ent.URL, ent.XPATH, key)
			// do not check again after 30 minutes
			ent.NextCheck = ent.NextCheck.Add(time.Minute * 30)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CancelWaitTime))
		defer cancel()
		if err := ent.Save(ctx, EntityType, DsClient, true); err != nil {
			log.Printf("ERROR: failed to save entity [%s] with %v", ent.Name, err)
		}
	}()

	msg := pushover.Message{
		Msg:  ent.String(),
		User: ent.Options.User,
	}

	var tracker trackers.Tracker = trackers.SimpleTracker
	if ent.Options.UseChrome {
		tracker = trackers.ChromeTracker
	}

	content, err := tracker(&ent.URL, &ent.XPATH)
	if err != nil {
		msg.Title = fmt.Sprintf("[%s] Alert: failed to fetch price `%v`!", ent.Name, err)
		pushClient.Send(&msg)
		return
	}
	if ent.History == nil {
		log.Println("WARN: zero price history")
		ent.History = []models.DataPoint{{Price: content, Timestamp: time.Now()}}
		return
	}

	last := ent.History[len(ent.History)-1]
	thisP, err := strconv.ParseFloat(priceRegex.FindString(content), 32)
	if err != nil {
		msg.Title = fmt.Sprintf("[%s] Alert: failed to convert price `%s`!", ent.Name, content)
		pushClient.Send(&msg)
		return
	}

	// update history & save entity
	ent.History = append(ent.History, models.DataPoint{Price: content, Timestamp: time.Now()})
	ent.NextCheck = time.Now().Add(time.Minute * time.Duration(ent.Options.CheckFreq))
	deltaRecordCnt := len(ent.History) - int(ent.Options.MaxRecords)
	if deltaRecordCnt > 0 {
		ent.History = ent.History[deltaRecordCnt:]
	}
	msg.Msg = ent.String() // update message
	// send alert
	if ent.Options.AlertType == "onChange" && content != last.Price {
		msg.Title = fmt.Sprintf("[%s] Alert: price changes to %s!", ent.Name, content)
		pushClient.Send(&msg)
	}
	if ent.Options.AlertType == "threshold" && ent.Options.Threshold >= float32(thisP) {
		msg.Title = fmt.Sprintf("[%s] Alert: price drops to %s!", ent.Name, content)
		pushClient.Send(&msg)
	}
	return
}

// Refresh refreshes prices from datastore
func Refresh(pushClient *pushover.Client, fetchLimit int) {
	log.Println("INFO: Refresh started")
	entities := FetchData(fetchLimit)
	for _, ent := range entities {
		if ent.K == nil {
			continue
		}
		log.Printf("INFO: processing [%s] XPATH (%s) at %s", ent.Name, ent.XPATH, ent.URL)
		processEntity(&ent, pushClient)
	}
	log.Println("INFO: Refresh ended")
}
