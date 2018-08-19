package gutils

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/xiahongze/pricetracker/types"
	"google.golang.org/api/iterator"
)

// FetchData fetches n records from datastore that needs to be checked
func FetchData(n int) []types.Entity {
	entities := make([]types.Entity, n, n)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	q := datastore.NewQuery(EntityType).Filter("NextCheck <", time.Now()).Limit(n)

	i := 0
	for t := DsClient.Run(ctx, q); ; {
		entity := &entities[i]
		_, err := t.Next(entity)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Handle error.
		}
		i++
	}
	return entities
}
