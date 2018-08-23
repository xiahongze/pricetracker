package models

import "cloud.google.com/go/datastore"

// CreateResponse returns whether create is successful
type CreateResponse struct {
	OK      bool           `json:"ok"`
	Message string         `json:"message"`
	Key     *datastore.Key `json:"key,omitempty"`
}

// ReadOrDelResponse returns whether read or del is successful
type ReadOrDelResponse struct {
	OK      bool    `json:"ok"`
	Message string  `json:"message"`
	Entity  *Entity `json:"entity,omitempty"`
}
