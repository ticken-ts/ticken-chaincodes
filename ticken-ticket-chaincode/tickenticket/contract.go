package tickenticket

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Issue(ctx TickenTxContext, payload *ticketPayload) (*Ticket, error) {
	ticketList := ctx.GetTicketList()
	tickenEventInvoker := ctx.GetTickenEventInvoker()

	payload.Sanitize()
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	eventID, ticketID, section, owner := payload.Deconstruct()

	ticketWithSameKey, err := ticketList.GetTicket(eventID, ticketID)
	if ticketWithSameKey != nil {
		return nil, fmt.Errorf("ticket %s already exists for event %s", ticketID, eventID)
	}

	_, err = tickenEventInvoker.GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	isAvailable, err := tickenEventInvoker.IsAvailable(eventID, section)
	if !isAvailable || err != nil {
		return nil, err
	}

	ticket := NewTicket(ticketID, eventID, section, owner)

	if err = ticketList.AddTicket(ticket); err != nil {
		return nil, err
	}
	if err = tickenEventInvoker.AddTicket(eventID, ticketID); err != nil {
		return nil, err
	}

	return ticket, nil
}
