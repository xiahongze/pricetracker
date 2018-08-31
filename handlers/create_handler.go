package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/trackers"
)

// CreateHandler handles create request
func CreateHandler(c echo.Context) error {
	req := &models.CreateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}
	content, ok := trackers.SimpleTracker(&req.URL, &req.XPATH)
	if !ok {
		return errors.New(content)
	}

	// add datastore handlers
	entity := gutils.ConvReq2Ent(req)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gutils.CancelWaitTime))
	defer cancel()
	err := entity.Save(ctx, gutils.EntityType, gutils.DsClient)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, entity)
}
