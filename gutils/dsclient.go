package gutils

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/datastore"
)

// ProjectName is the project name in your GCP
var ProjectName = "project-order-management"

// EntityType is the EntityType in your datastore
var EntityType = "price-tracks"

// DsClient is the global datastore client
var DsClient *datastore.Client

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
		log.Fatal("failed to new a DsClient", err)
	}
}
