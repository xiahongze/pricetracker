package gutils

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

// ProjectName is the project name in your GCP
var ProjectName = "project-order-management"

// EntityType is the EntityType in your datastore
var EntityType = "price-tracks"

// DsClient is the global datastore client
var DsClient *datastore.Client

// CancelWaitTime is the default GCP wait time for an operation
var CancelWaitTime = time.Second * 30

func init() {
	if v, ok := os.LookupEnv("PROJECT_NAME"); ok {
		ProjectName = v
	}
	if v, ok := os.LookupEnv("ENTITY_TYPE"); ok {
		EntityType = v
	}
	ctx := context.Background()
	var err error
	DsClient, err = datastore.NewClient(ctx, ProjectName)
	if err != nil {
		log.Fatal("ERROR: failed to new a DsClient", err)
	}
}
