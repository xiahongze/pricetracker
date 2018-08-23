package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/trackers"
)

func writeCreateResponse(w http.ResponseWriter, msg string, status int, OK bool, Key *datastore.Key) {
	resp := models.CreateResponse{OK: OK, Message: msg, Key: Key}
	b, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(b)
}

func create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	req := models.CreateRequest{}
	json.Unmarshal(body, &req)
	if content, ok := req.Validate(); !ok {
		resp := models.CreateResponse{OK: false, Message: content}
		b, _ := json.Marshal(resp)
		writeCreateResponse(w, string(b), http.StatusBadRequest, false, nil)
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
	entity := gutils.ConvReq2Ent(req)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	err := entity.Save(ctx, gutils.EntityType, gutils.DsClient)

	if err != nil {
		writeCreateResponse(w, err.Error(), http.StatusInternalServerError, false, nil)
		return
	}

	writeCreateResponse(w, "success", http.StatusOK, true, entity.K)
}

// CreateHandler returns the value if any from the xpath of the url
var CreateHandler = http.HandlerFunc(create)
