package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/models"
	"github.com/xiahongze/pricetracker/pushover"
	"github.com/xiahongze/pricetracker/trackers"
)

// MakeCreate creates create handler request
func MakeCreate(client *pushover.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &models.CreateRequest{}
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if msg, ok := req.Validate(); !ok {
			return c.String(http.StatusBadRequest, msg)
		}

		var (
			content   string
			ok        bool
			useChrome = true
		)

		if req.Options.UseChrome == nil || !*req.Options.UseChrome {
			content, ok = trackers.SimpleTracker(&req.URL, &req.XPATH)
		}
		if !ok {
			req.Options.UseChrome = &useChrome
			log.Println("INFO: Resorting to Chrome")
		}
		if req.Options.UseChrome != nil && *req.Options.UseChrome {
			if content, ok = trackers.ChromeTracker(&req.URL, &req.XPATH); !ok {
				return c.String(http.StatusBadRequest, content)
			}
		}

		// check content as expected
		if content != req.ExpectedPrice {
			return c.String(http.StatusExpectationFailed,
				fmt.Sprintf("expected price (%s) != extracted price (%s)",
					req.ExpectedPrice, content))
		}

		// add datastore handlers
		entity := gutils.ConvReq2Ent(req)
		ctx, cancel := context.WithTimeout(context.Background(), gutils.CancelWaitTime)
		defer cancel()
		if err := entity.Save(ctx, gutils.EntityType, gutils.DsClient, false); err != nil {
			return err
		}

		msg := pushover.Message{
			Title: fmt.Sprintf("[%s] INFO: %s!", entity.Name, "Created"),
			Msg:   entity.String(),
		}
		client.Send(&msg)

		return c.JSON(http.StatusCreated, entity)
	}
}
