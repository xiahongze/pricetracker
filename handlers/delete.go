package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/email"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
)

// Delete deletes the record and returns the given entity record
func Delete(c echo.Context) error {
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

	ctx1, cancel1 := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
	defer cancel1()
	if err := gutils.DsClient.Delete(ctx1, req.Key); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	subject := fmt.Sprintf("[%s] <%s> INFO: %s!", email.Identity, entity.Name, "Deleted")
	entity.SendEmail(&subject)

	return c.JSON(http.StatusOK, entity)
}
