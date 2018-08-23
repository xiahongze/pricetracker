package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
)

func writeReadResponse(w http.ResponseWriter, status int, OK bool, msg string, Entity *models.Entity) {
	resp := models.ReadOrDelResponse{OK: OK, Entity: Entity, Message: msg}
	b, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(b)
}

func read(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	req := models.ReadOrDelRequest{}
	json.Unmarshal(body, &req)
	if content, ok := req.Validate(); !ok {
		writeReadResponse(w, http.StatusBadRequest, false, content, nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	entity := &models.Entity{}
	err := gutils.DsClient.Get(ctx, req.Key, entity)
	if err != nil {
		writeReadResponse(w, http.StatusNotFound, false, err.Error(), nil)
		return
	}

	writeReadResponse(w, http.StatusOK, true, "", entity)
}

// ReadHandler returns the given entity record
var ReadHandler = http.HandlerFunc(read)
