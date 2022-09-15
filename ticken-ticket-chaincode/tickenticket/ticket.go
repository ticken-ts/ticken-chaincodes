package tickenticket

import (
	"encoding/json"
	"fmt"
	"strings"
	"ticken-ticket-contract/ledgerapi"
)

type Status string

const (
	ISSUED  Status = "Issued"
	SIGNED  Status = "Signed"
	SCANNED Status = "Scanned"
)

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Section  string `json:"section"`
	Owner    string `json:"owner"`
	Status   Status `json:"status"`
}

func TicketCreateKey(eventID string, ticketID string) string {
	return ledgerapi.MakeKey(eventID, ticketID)
}

func TicketDeserialize(jsonBytes []byte, ticket *Ticket) error {
	err := json.Unmarshal(jsonBytes, ticket)

	if err != nil {
		return fmt.Errorf("error deserializing ticken-event. %s", err.Error())
	}

	return nil
}

// ******************** Primitives ********************** //

func NewTicket(ticketID string, eventID string, section string, owner string) (*Ticket, error) {
	ticket := new(Ticket)

	ticket.Status = ISSUED

	ticket.Owner = strings.TrimSpace(owner)
	ticket.EventID = strings.TrimSpace(eventID)
	ticket.TicketID = strings.TrimSpace(ticketID)
	ticket.Section = strings.ToUpper(strings.TrimSpace(section))

	if err := ticket.validate(); err != nil {
		return nil, err
	}

	return ticket, nil
}

// The following implementations are required
// to implement the State.

func (ticket *Ticket) GetKey() string {
	return ledgerapi.MakeKey(ticket.EventID, ticket.TicketID)
}

func (ticket *Ticket) Serialize() ([]byte, error) {
	return json.Marshal(ticket)
}

// ****************************************************** //

func (ticket *Ticket) validate() error {
	// TODO -> check if necessary to validate  uuid format for ID's

	if len(ticket.Owner) == 0 {
		return fmt.Errorf("owner is mandatory")
	}

	if len(ticket.TicketID) == 0 {
		return fmt.Errorf("ticketID is mandatory")
	}

	if len(ticket.EventID) == 0 {
		return fmt.Errorf("eventID is mandatory")
	}

	if len(ticket.Section) == 0 {
		return fmt.Errorf("section is mandatory")
	}

	return nil
}
