package tickenticket

import (
	"encoding/json"
	"fmt"
	"strings"

	ledgerapi "ticken-ticket-contract/ledgerapi"
)

// ********************** Ticket ************************ //

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

// ****************************************************** //

// ******************** Primitives ********************** //

func NewTicket(ticketID string, eventID string, section string, owner string) *Ticket {
	ticket := new(Ticket)

	ticket.Status = ISSUED
	ticket.Owner = owner
	ticket.Section = strings.ToUpper(section)
	ticket.EventID = eventID
	ticket.TicketID = ticketID

	return ticket
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
