package utils

import (
	"time"

	"github.com/xiahongze/pricetracker/types"
)

// ConvReq2Ent converts CreateRequest to Entity in datastore
func ConvReq2Ent(req types.CreateRequest) types.Entity {
	return types.Entity{
		Options:   *req.Options,
		URL:       req.URL,
		XPATH:     req.XPATH,
		NextCheck: time.Now().Add(time.Minute * time.Duration(req.Options.CheckFreq)),
		History:   []types.DataPoint{types.DataPoint{Timestamp: time.Now(), Price: req.ExpectedPrice}},
	}
}
