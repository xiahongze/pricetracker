package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/pushover"
)

// MakeUpdate creates the Update handler
func MakeUpdate(client *pushover.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &models.UpdateRequest{}
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
		defer cancel()
		entity := &models.Entity{}
		if err := gutils.DsClient.Get(ctx, req.Key, entity); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// update relevant fields
		if req.Options.AlertType != "" {
			entity.Options.AlertType = req.Options.AlertType
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
		if req.UseChrome != nil {
			entity.Options.UseChrome = *req.UseChrome
		}

		if req.Name != "" {
			entity.Name = req.Name
		}
		if req.URL != "" {
			entity.URL = req.URL
		}
		if req.XPATH != "" {
			entity.XPATH = req.XPATH
		}

		// update entity
		ctx1, cancel1 := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
		defer cancel1()
		if err := entity.Save(ctx1, gutils.EntityType, gutils.DsClient, false); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		msg := pushover.Message{
			Title: fmt.Sprintf("[%s] INFO: %s!", entity.Name, "Updated"),
			Msg:   entity.String(),
		}
		client.Send(&msg)

		return c.JSON(http.StatusOK, entity)
	}
}
