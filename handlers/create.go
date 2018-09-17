package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/trackers"
)

// Create handles create request
func Create(c echo.Context) error {
	req := &models.CreateRequest{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if msg, ok := req.Validate(); !ok {
		return c.String(http.StatusBadRequest, msg)
	}

	if content, ok := trackers.SimpleTracker(&req.URL, &req.XPATH); !ok {
		return c.String(http.StatusBadRequest, content)
	}

	// add datastore handlers
	entity := gutils.ConvReq2Ent(req)
	ctx, cancel := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel()
	if err := entity.Save(ctx, gutils.EntityType, gutils.DsClient); err != nil {
		return err
	}

	sendEmail(&entity, "Created")

	return c.JSON(http.StatusCreated, entity)
}
