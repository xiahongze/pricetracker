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

func processEntity(ent *models.Entity, pushClient *pushover.Client) {
	// save the entity before returning
	defer func() {
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
	if ent.Options.UseChrome != nil && *ent.Options.UseChrome {
		tracker = trackers.ChromeTracker
	}

	content, ok := tracker(&ent.URL, &ent.XPATH)
	if !ok {
		log.Println("ERROR: failed to fetch price.", content)
		key, _ := ent.K.MarshalJSON()
		log.Printf("URL: %s\nXPATH: %s\nKey: %s", ent.URL, ent.XPATH, key)
		msg.Title = fmt.Sprintf("[%s] Alert: failed to fetch price because`%s`!", ent.Name, content)
		pushClient.Send(&msg)
		// do not check again after 30 minutes
		ent.NextCheck = ent.NextCheck.Add(time.Minute * 30)
		return
	}
	if ent.History == nil {
		log.Println("WARN: zero price history.", ent)
		ent.History = []models.DataPoint{{Price: content, Timestamp: time.Now()}}
		return
	}

	last := ent.History[len(ent.History)-1]
	thisP, err := strconv.ParseFloat(priceRegex.FindString(content), 32)
	if err != nil {
		log.Println("ERROR: failed to convert price", err, "this price:", content)
		msg.Title = fmt.Sprintf("[%s] Alert: failed to convert price `%s`!", ent.Name, content)
		pushClient.Send(&msg)
		// do not check again after 30 minutes
		ent.NextCheck = ent.NextCheck.Add(time.Minute * 30)
		return
	}

	// update history & save entity
	ent.History = append(ent.History, models.DataPoint{Price: content, Timestamp: time.Now()})
	ent.NextCheck = time.Now().Add(time.Minute * time.Duration(ent.Options.CheckFreq))
	deltaRecordCnt := len(ent.History) - int(ent.Options.MaxRecords)
	if deltaRecordCnt > 0 {
		ent.History = ent.History[deltaRecordCnt:]
	}
	// send alert
	if ent.Options.AlertType == "onChange" && content != last.Price {
		msg.Title = fmt.Sprintf("[%s] Alert: price changes to %s!", ent.Name, content)
		pushClient.Send(&msg)
	}
	if ent.Options.AlertType == "threshold" && ent.Options.Threshold >= float32(thisP) {
		msg.Title = fmt.Sprintf("[%s] Alert: price drops to %s!", ent.Name, content)
		pushClient.Send(&msg)
	}
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
