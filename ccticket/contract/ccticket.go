package contract

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/ticken-ts/ticken-chaincodes/common"
)

type Contract struct {
	contractapi.Contract
}

const Name = "cc-ticket"

// ********** cc-event integration ********** //

const ccEventName = "cc-event"
const ccEventSellTicketFunc = "SellTicket"

// *****+************************************ //

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

func (c *Contract) Issue(ctx common.ITickenTxContext, ticketID, eventID, section, ownerID string) (*Ticket, error) {
	_, err := c.GetTicket(ctx, ticketID)
	if err != nil {
		return nil, err // this error is already formatted
	}

	ownerIDParsed, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, ccErr("error parsing owner id: %v", err)
	}
	eventIDParsed, err := uuid.Parse(eventID)
	if err != nil {
		return nil, ccErr("error parsing event id: %v", err)
	}
	ticketIDParsed, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, ccErr("error parsing ticket id: %v", err)
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
		return nil, ccErr("failed to serialize ticket: %v", err)
	}

	// add ticket into the chaincode cc-event
	// note: this operation is atomically handled
	// by the orderers. So, the ticket and the ticket
	// count are updated simultaneously in the same tx
	_, err = ctx.GetInvoker(ccEventName).Invoke(ccEventSellTicketFunc, eventID, section)
	if err != nil {
		return nil, ccErr(err.Error())
	}

	if err := ctx.GetStub().PutState(ticketID, ticketJSON); err != nil {
		return nil, ccErr("failed to updated the state: %v", err)
	}

	return &ticket, nil
}

func (c *Contract) GetTicket(ctx contractapi.TransactionContextInterface, ticketID string) (*Ticket, error) {
	ticketJSON, err := ctx.GetStub().GetState(ticketID)
	if err != nil {
		return nil, ccErr("failed to read ticket: %v", err)
	}
	if ticketJSON == nil {
		return nil, ccErr("the event %s does not exist", ticketID)
	}

	var ticket Ticket
	if err := json.Unmarshal(ticketJSON, &ticket); err != nil {
		return nil, ccErr("failed to deserialize ticket: %v", err)
	}

	return &ticket, err
}

func ccErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args)
	return fmt.Errorf("[%s] | %v", Name, msg)
}
