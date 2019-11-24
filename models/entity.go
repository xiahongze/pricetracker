package models

import (
	"context"
	"encoding/json"
	"errors"
	"log"
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

// String returns the JSON String representation
func (entity *Entity) String() string {
	// marshal the entity as the message
	b, err := json.MarshalIndent(entity, "", "    ")
	if err != nil {
		log.Print("failed to marshal entity", err)
		return ""
	}
	return string(b)
}
