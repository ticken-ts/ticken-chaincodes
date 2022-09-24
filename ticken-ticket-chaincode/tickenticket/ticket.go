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
	TicketID  string `json:"ticket_id"`
	EventID   string `json:"event_id"`
	Section   string `json:"section"`
	Owner     string `json:"owner"`
	Status    Status `json:"status"`
	Signature string `json:"signature"`
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

func (ticket *Ticket) IsFrom(owner string) bool {
	return ticket.Owner == strings.TrimSpace(owner)
}

func (ticket *Ticket) Sign(signer string, signature string) error {
	// this method will contain the login to sing
	// itself. Maybe we should receive a private key
	// or something more secure to sign it
	if !ticket.isAllowToSign() {
		return fmt.Errorf("ticket is not allow to sign")
	}

	if !ticket.IsFrom(signer) {
		return fmt.Errorf("ticket is not from signer %s", signer)
	}

	if len(signature) == 0 {
		return fmt.Errorf("ticket signature can't be empty")
	}

	ticket.Signature = signature
	ticket.Status = SIGNED
	return nil
}

func (ticket *Ticket) Scan() error {
	if ticket.Status == SCANNED {
		return fmt.Errorf("ticket is already scanned")
	}

	if ticket.Status == ISSUED {
		return fmt.Errorf("ticket is not signed")
	}

	ticket.Status = SCANNED
	return nil
}

func (ticket *Ticket) isAllowToSign() bool {
	return ticket.Status == ISSUED
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
	// TODO -> check if necessary to validate uuid format for ID's

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
