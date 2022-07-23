package ticken_ticket

import (
	"encoding/json"
	"fmt"
	ledgerapi "ticken-ticket-contract/ledger-api"
)

// ********************** Ticket ************************ //

type Status uint

const (
	ISSUED Status = iota + 1
	SIGNED
	SCANNED
)

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Owner    string `json:"owner"`
	Status   Status `json:"status"`
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

func (status Status) String() string {
	statusNames := []string{"ISSUED", "SIGNED", "SCANNED"}

	if status < ISSUED || status > SIGNED {
		return "UNKNOWN"
	}

	return statusNames[status-1]
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
