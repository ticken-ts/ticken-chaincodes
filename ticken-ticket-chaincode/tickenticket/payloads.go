package tickenticket

import (
	"fmt"
	"strings"
)

type ticketPayload struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Section  string `json:"section"`
	Owner    string `json:"owner"`
}

func (payload *ticketPayload) Sanitize() {
	payload.TicketID = strings.TrimSpace(payload.TicketID)
	payload.EventID = strings.TrimSpace(payload.EventID)
	payload.Section = strings.TrimSpace(payload.Section)
	payload.Owner = strings.TrimSpace(payload.Owner)
}

func (payload *ticketPayload) Validate() error {
	if len(payload.TicketID) == 0 {
		return fmt.Errorf("ticket_id is mandatory")
	}
	return nil
}

func (payload *ticketPayload) Deconstruct() (string, string, string, string) {
	return payload.EventID, payload.TicketID, payload.Section, payload.Owner
}
