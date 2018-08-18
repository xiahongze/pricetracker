package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/xiahongze/pricetracker/trackers"
	"github.com/xiahongze/pricetracker/types"
	"github.com/xiahongze/pricetracker/types/utils"
)

func writeCreateResponse(w http.ResponseWriter, msg string, status int, OK bool, Key *datastore.Key) {
	resp := types.CreateResponse{OK: OK, Message: msg, Key: Key}
	b, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(b)
}

func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	req := types.CreateRequest{}
	json.Unmarshal(body, &req)
	if content, ok := req.Validate(); !ok {
		resp := types.CreateResponse{OK: false, Message: content}
		b, _ := json.Marshal(resp)
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}

	content, ok := trackers.SimpleTracker(&req.URL, &req.XPATH)
	if !ok {
		writeCreateResponse(w, content, http.StatusBadRequest, false, nil)
		return
	}

	if content != req.ExpectedPrice {
		msg := "expectedPrice did not match with the fetched result (" + content + ")"
		writeCreateResponse(w, msg, http.StatusExpectationFailed, false, nil)
		return
	}

	// add datastore handlers
	entity := utils.ConvReq2Ent(req)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	err := entity.Save(ctx, EntityType, dsClient)

	if err != nil {
		writeCreateResponse(w, err.Error(), http.StatusInternalServerError, false, nil)
		return
	}

	writeCreateResponse(w, "success", http.StatusOK, true, entity.K)
}

// CreateHandler returns the value if any from the xpath of the url
var CreateHandler = http.HandlerFunc(create)
