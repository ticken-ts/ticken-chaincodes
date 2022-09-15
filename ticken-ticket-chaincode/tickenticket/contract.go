package tickenticket

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Issue(ctx TickenTxContext, ticketID string, eventID string, section string, owner string) (*Ticket, error) {
	ticketList := ctx.GetTicketList()
	tickenEventInvoker := ctx.GetTickenEventInvoker()

	ticket, err := NewTicket(ticketID, eventID, section, owner)
	if err != nil {
		return nil, err
	}

	ticketIsDuplicated, err := ticketList.TicketExist(eventID, ticketID)
	if ticketIsDuplicated || err != nil {
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("ticket %s already exists for event %s", ticketID, eventID)
		}
	}

	eventExists, err := tickenEventInvoker.EventExists(eventID)
	if !eventExists || err != nil {
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("eventID %s doesn't exists", eventID)
		}
	}

	isAvailable, err := tickenEventInvoker.IsAvailable(eventID, section)
	if !isAvailable || err != nil {
		return nil, err
	}

	if err = ticketList.AddTicket(ticket); err != nil {
		return nil, err
	}
	if err = tickenEventInvoker.AddTicket(eventID, ticketID); err != nil {
		return nil, err
	}

	return ticket, nil
}
