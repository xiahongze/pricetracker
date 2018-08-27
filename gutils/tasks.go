package gutils

import (
	"context"
	"fmt"
	"log"
	"math"
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

func composeEmail(ent models.Entity) string {
	key, err := ent.K.MarshalJSON()
	if err != nil {
		key = []byte("unrecognizable key")
	}
	return fmt.Sprintf(`Link to the website: %s
Name: %s
XPATH: %s
Key: %s`, ent.URL, ent.Name, ent.XPATH, key)
}

// Refresh refreshes prices from datastore
func Refresh() {
	entities := FetchData(fetchLimit)
	if l := len(entities); l > 0 {
		log.Println("INFO: have fetched", l, "entities")
		for _, ent := range entities {
			if content, ok := trackers.SimpleTracker(&ent.URL, &ent.XPATH); ok {
				if ent.History != nil {
					last := ent.History[len(ent.History)-1]
					lastP, err := strconv.ParseFloat(priceRegex.FindString(last.Price), 32)
					if err != nil {
						log.Println("ERROR: failed to convert price", err, "last price:", last.Price)
						return
					}
					thisP, err := strconv.ParseFloat(priceRegex.FindString(content), 32)
					if err != nil {
						log.Println("ERROR: failed to convert price", err, "this price:", content)
						return
					}
					if ent.Options.AlertType == "onChange" && math.Abs(lastP-thisP) > 1e-3 {
						subjec := fmt.Sprintf("[%s] <%s> Alert: price changes to %s!", email.Identity, ent.Name, content)
						email.Send(composeEmail(ent), subjec, ent.Options.Email)
					} else if ent.Options.AlertType == "threshold" && ent.Options.Threshold >= float32(thisP) {
						subjec := fmt.Sprintf("[%s] <%s> Alert: price drops to %s!", email.Identity, ent.Name, content)
						email.Send(composeEmail(ent), subjec, ent.Options.Email)
					}

					// update history & save entity
					ent.History = append(ent.History, models.DataPoint{Price: content, Timestamp: time.Now()})
					ent.NextCheck = time.Now().Add(time.Minute * time.Duration(ent.Options.CheckFreq))
					if len(ent.History) > int(ent.Options.MaxRecords) {
						ent.History = ent.History[:ent.Options.MaxRecords]
					}
				} else {
					log.Println("WARN: zero price history.", ent)
					ent.History = []models.DataPoint{models.DataPoint{Price: content, Timestamp: time.Now()}}
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CancelWaitTime))
				defer cancel()
				if err := ent.Save(ctx, EntityType, DsClient); err != nil {
					log.Println("ERROR: failed to save entity:", err, ". Entity: ", ent)
				}
			} else {
				log.Println("ERROR: failed to fetch price.", content)
			}
		}
	} else {
		log.Println("INFO: No updates")
	}
}
