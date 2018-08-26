package gutils

import (
	"time"

	"github.com/xiahongze/pricetracker/models"
)

// ConvReq2Ent converts CreateRequest to Entity in datastore
func ConvReq2Ent(req models.CreateRequest) models.Entity {
	return models.Entity{
		Options:   *req.Options,
		URL:       req.URL,
		XPATH:     req.XPATH,
		NextCheck: time.Now().Add(time.Minute * time.Duration(req.Options.CheckFreq)),
		History:   []models.DataPoint{models.DataPoint{Timestamp: time.Now(), Price: req.ExpectedPrice}},
	}
}