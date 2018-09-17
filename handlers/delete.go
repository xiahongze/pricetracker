package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
)

// Delete deletes the record and returns the given entity record
func Delete(c echo.Context) error {
	req := &models.ReadOrDelRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel()
	entity := &models.Entity{}
	if err := gutils.DsClient.Get(ctx, req.Key, entity); err != nil {
		return err
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel1()
	if err := gutils.DsClient.Delete(ctx1, req.Key); err != nil {
		return err
	}

	sendEmail(entity, "Deleted")

	return c.JSON(http.StatusOK, entity)
}
