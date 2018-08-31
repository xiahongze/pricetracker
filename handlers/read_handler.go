package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
)

// ReadHandler returns the given entity record
func ReadHandler(c echo.Context) error {
	req := &models.ReadOrDelRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	entity := &models.Entity{}
	if err := gutils.DsClient.Get(ctx, req.Key, entity); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, entity)
}
