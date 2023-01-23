package contract

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

const Name = "cc-ticket"

type Ticket struct {
	TicketID uuid.UUID `json:"ticket_id"`
	Status   Status    `json:"status"`

	EventID uuid.UUID `json:"event_id"`
	Section string    `json:"section"`

	// represents the owner id in the
	// web service database
	OwnerID uuid.UUID `json:"owner"`
}

type Status string

const (
	ISSUED Status = "issued"
)

func (c *Contract) Issue(ctx contractapi.TransactionContextInterface, ticketID, eventID, section, ownerID string) (*Ticket, error) {
	_, err := c.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	ownerIDParsed, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, ccError(fmt.Errorf("error parsing owner id: %v", err))
	}
	eventIDParsed, err := uuid.Parse(eventID)
	if err != nil {
		return nil, ccError(fmt.Errorf("error parsing event id: %v", err))
	}
	ticketIDParsed, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, ccError(fmt.Errorf("error parsing ticket id: %v", err))
	}

	ticket := Ticket{
		TicketID: ticketIDParsed,

		EventID: eventIDParsed,
		Section: section,
		Status:  ISSUED,

		OwnerID: ownerIDParsed,
	}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		return nil, ccError(fmt.Errorf("failed to serialize ticket: %v", err))
	}

	if err := ctx.GetStub().PutState(ticketID, ticketJSON); err != nil {
		return nil, ccError(fmt.Errorf("failed to updated the state: %v", err))
	}

	return &ticket, nil
}

func (c *Contract) GetTicket(ctx contractapi.TransactionContextInterface, ticketID string) (*Ticket, error) {
	ticketJSON, err := ctx.GetStub().GetState(ticketID)
	if err != nil {
		return nil, ccError(fmt.Errorf("failed to read ticket: %v", err))
	}
	if ticketJSON == nil {
		return nil, ccError(fmt.Errorf("the event %s does not exist", ticketID))
	}

	var ticket Ticket
	err = json.Unmarshal(ticketJSON, &ticket)
	if err != nil {
		return nil, ccError(fmt.Errorf("failed to descerialize ticket: %v", err))
	}

	return &ticket, nil
}

func ccError(err error) error {
	return fmt.Errorf("[%s] | %v", Name, err)
}
