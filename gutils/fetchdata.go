package gutils

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/xiahongze/pricetracker/models"
	"google.golang.org/api/iterator"
)

// FetchData fetches n records from datastore that needs to be checked
func FetchData(n int) []models.Entity {
	entities := make([]models.Entity, n, n)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CancelWaitTime))
	defer cancel()

	q := datastore.NewQuery(EntityType).Filter("NextCheck <", time.Now()).Limit(n)

	i := 0
	for t := DsClient.Run(ctx, q); i < n; i++ {
		entity := &entities[i]
		_, err := t.Next(entity)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error.
			log.Println("ERROR:", err)
		}
	}
	return entities
}
