package gutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/xiahongze/pricetracker/email"
	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/trackers"
)

var fetchLimit = 10
var priceRegex, _ = regexp.Compile("\\d+\\.?\\d{0,}")

func init() {
	if v, ok := os.LookupEnv("FETCH_LIMIT"); ok {
		tmpI, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalln("ERROR: ", err)
		}
		fetchLimit = tmpI
	}
}

func processEntity(ent *models.Entity) {
	if ent.K == nil {
		return
	}
	// save the entity before returning
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CancelWaitTime))
	defer func() {
		if err := ent.Save(ctx, EntityType, DsClient); err != nil {
			log.Println("ERROR: failed to save entity:", err, ". Entity: ", ent)
		}
		cancel()
	}()

	var tracker trackers.Tracker
	switch ent.Options.UseChrome {
	case true:
		tracker = trackers.ChromeTracker
	default:
		tracker = trackers.SimpleTracker
	}

	content, ok := tracker(&ent.URL, &ent.XPATH)
	if !ok {
		log.Println("ERROR: failed to fetch price.", content)
		key, _ := ent.K.MarshalJSON()
		log.Printf("URL: %s\nXPATH: %s\nKey: %s", ent.URL, ent.XPATH, key)
		subject := fmt.Sprintf("[%s] <%s> Alert: failed to fetch price for reason `%s`!", email.Identity, ent.Name, content)
		ent.SendEmail(&subject)
		return
	}
	if ent.History == nil {
		log.Println("WARN: zero price history.", ent)
		ent.History = []models.DataPoint{models.DataPoint{Price: content, Timestamp: time.Now()}}
		return
	}

	last := ent.History[len(ent.History)-1]
	thisP, err := strconv.ParseFloat(priceRegex.FindString(content), 32)
	if err != nil {
		log.Println("ERROR: failed to convert price", err, "this price:", content)
		return
	}

	// update history & save entity
	ent.History = append(ent.History, models.DataPoint{Price: content, Timestamp: time.Now()})
	ent.NextCheck = time.Now().Add(time.Minute * time.Duration(ent.Options.CheckFreq))
	if len(ent.History) > int(ent.Options.MaxRecords) {
		ent.History = ent.History[:ent.Options.MaxRecords]
	}
	// send alert
	if ent.Options.AlertType == "onChange" && content != last.Price {
		subject := fmt.Sprintf("[%s] <%s> Alert: price changes to %s!", email.Identity, ent.Name, content)
		ent.SendEmail(&subject)
	}
	if ent.Options.AlertType == "threshold" && ent.Options.Threshold >= float32(thisP) {
		subject := fmt.Sprintf("[%s] <%s> Alert: price drops to %s!", email.Identity, ent.Name, content)
		ent.SendEmail(&subject)
	}
}

// Refresh refreshes prices from datastore
func Refresh() {
	log.Println("INFO: Refresh started")
	entities := FetchData(fetchLimit)
	for _, ent := range entities {
		processEntity(&ent)
	}
	log.Println("INFO: Refresh ended")
}
