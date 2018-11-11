package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
)

// Read returns the given entity record
func Read(c echo.Context) error {
	req := &models.ReadOrDelRequest{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel()
	entity := &models.Entity{}
	if err := gutils.DsClient.Get(ctx, req.Key, entity); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, entity)
}
