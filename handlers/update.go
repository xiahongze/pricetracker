package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
)

// Update updates the entry given the key
func Update(c echo.Context) error {
	req := &models.UpdateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel()
	entity := &models.Entity{}
	if err := gutils.DsClient.Get(ctx, req.Key, entity); err != nil {
		return err
	}

	// update relevant fields
	if req.Options.AlertType != "" {
		entity.Options.AlertType = req.Options.AlertType
	}
	if req.Options.Email != "" {
		entity.Options.Email = req.Options.Email
	}
	if req.Options.CheckFreq != 0 {
		entity.Options.CheckFreq = req.Options.CheckFreq
	}
	if req.Options.MaxRecords != 0 {
		entity.Options.MaxRecords = req.Options.MaxRecords
	}
	if req.Options.Threshold != 0 {
		entity.Options.Threshold = req.Options.Threshold
	}
	if req.Options.UseChrome {
		entity.Options.UseChrome = req.Options.UseChrome
	}

	// update entity
	ctx1, cancel1 := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel1()
	if err := entity.Save(ctx1, gutils.EntityType, gutils.DsClient); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, entity)
}
