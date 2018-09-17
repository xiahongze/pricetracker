package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xiahongze/pricetracker/email"
	"github.com/xiahongze/pricetracker/models"
)

func sendEmail(entity *models.Entity, info string) {
	subject := fmt.Sprintf("[%s] <%s> INFO: %s!", email.Identity, entity.Name, info)
	if b, err := json.MarshalIndent(entity, "", "    "); err == nil {
		if err := email.Send(string(b), subject, entity.Options.Email); err != nil {
			log.Print("failed to send email", err)
		}
	} else {
		log.Print("failed to marshal entity", err)
	}
}
