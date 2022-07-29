package ticken_ticket

import (
	"encoding/json"
	"fmt"
	ledgerapi "ticken-ticket-contract/ledger-api"
)

// ********************** Ticket ************************ //

type Status string

const (
	ISSUED  Status = "Issued"
	SIGNED  Status = "Signed"
	SCANNED Status = "Scanned"
)

type Ticket struct {
	TicketID   string `json:"ticket_id"`
	EventID    string `json:"event_id"`
	Owner      string `json:"owner"`
	Status     Status `json:"status"`
	Section    string `json:"section"`
	Subsection string `json:"subsection"`
}

func CreateTicketKey(eventID string, ticketID string) string {
	return ledgerapi.MakeKey(eventID, ticketID)
}

func Deserialize(jsonBytes []byte, ticket *Ticket) error {
	err := json.Unmarshal(jsonBytes, ticket)

	if err != nil {
		return fmt.Errorf("error deserializing ticken-event. %s", err.Error())
	}

	return nil
}

// ****************************************************** //

// ******************** Primitives ********************** //

func (ticket *Ticket) Init(ticketID string, eventID string, owner string) {
	ticket.Status = ISSUED
	ticket.Owner = owner
	ticket.EventID = eventID
	ticket.TicketID = ticketID
}

func (ticket *Ticket) IsIssued() bool {
	return ticket.Status == ISSUED
}

func (ticket *Ticket) IsSigned() bool {
	return ticket.Status == SIGNED
}

func (ticket *Ticket) IsScanned() bool {
	return ticket.Status == SCANNED
}

func (ticket *Ticket) Sign() {
	ticket.Status = SIGNED
}

func (ticket *Ticket) Scan() {
	ticket.Status = SCANNED
}

// The following implementations are required
// to implement the StateInterface.

func (ticket *Ticket) GetSplitKey() []string {
	return []string{ticket.TicketID, ticket.EventID}
}

func (ticket *Ticket) Serialize() ([]byte, error) {
	return json.Marshal(ticket)
}

// ****************************************************** //
