package models

import (
	"context"
	"encoding/json"
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
func (entity *Entity) Save(ctx context.Context, entTypName string, dsClient *datastore.Client) (err error) {
	defer func() {
		k, _ := json.Marshal(entity.K)
		if err != nil {
			log.Printf("ERROR: Entity: failed to save entity with %s\n", err)
		} else {
			log.Printf("INFO: Entity: saved K=%s\n", k)
		}
	}()

	if entity.K == nil {
		k := datastore.IncompleteKey(entTypName, nil)
		var key *datastore.Key
		key, err = dsClient.Put(ctx, k, entity)
		if err != nil {
			return err
		}
		entity.K = key
	} else {
		_, err = dsClient.Put(ctx, entity.K, entity)
		if err != nil {
			return err
		}
	}
	return nil
}
