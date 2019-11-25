package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
)

type (
	// DataPoint aka each data record
	DataPoint struct {
		Timestamp time.Time
		Price     string
	}

	// Entity is the data structure for datastore entry
	Entity struct {
		K         *datastore.Key `json:",omitempty" datastore:"__key__"`
		Options   Options        `json:",omitempty" datastore:",noindex"`
		URL       string         `json:",omitempty" datastore:",noindex"`
		XPATH     string         `json:",omitempty" datastore:",noindex"`
		Name      string         `json:",omitempty"`
		NextCheck time.Time      `json:",omitempty"`
		History   []DataPoint    `json:",omitempty" datastore:",noindex"`
	}
)

// Save saves the entry in the datastore
func (entity *Entity) Save(ctx context.Context, entTypName string, dsClient *datastore.Client, check bool) (err error) {
	defer func() {
		k, _ := json.Marshal(entity.K)
		if err != nil {
			log.Printf("ERROR: failed to save entity (K=%s) with %s\n", k, err)
			return
		}
		log.Printf("INFO: saved K=%s\n", k)
	}()

	if entity.K != nil {
		_, err = dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
			_ent := Entity{}
			if err := tx.Get(entity.K, &_ent); err != nil {
				return err
			}
			if _ent.NextCheck.After(time.Now()) && check {
				return errors.New("Entity was updated by another goroutine")
			}
			if _, err := tx.Mutate(datastore.NewUpdate(entity.K, entity)); err != nil {
				return err
			}
			return nil
		})
		return
	}

	entity.K = datastore.IncompleteKey(entTypName, nil)
	entity.K, err = dsClient.Put(ctx, entity.K, entity)

	return
}

// String returns a String representation
// History has been skipped to save text space but the last one will be there
func (entity *Entity) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Product [%s] - K=%s\nURL: %s\nXPATH: %s\nHistory:\n",
		entity.Name, entity.K.Encode(), entity.URL, entity.XPATH))
	lastPrice := ""
	for i, data := range entity.History {
		if lastPrice != data.Price || i == len(entity.History)-1 {
			sb.WriteString(data.Timestamp.Format(time.RFC822))
			sb.WriteString("\t")
			sb.WriteString(data.Price)
			sb.WriteString("\n")
			lastPrice = data.Price
		}
	}
	return sb.String()
}
