package tickenticket

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Issue(ctx TransactionContextInterface, payload *ticketPayload) (*Ticket, error) {
	ticketList := ctx.GetTicketList()
	tickenEventCCInvoker := ctx.GetTickenEventCCInvoker()

	payload.Sanitize()
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	eventID, ticketID, section, owner := payload.Deconstruct()

	event, err := tickenEventCCInvoker.GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	txCreationTimestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return nil, err
	}
	if !event.TicketSellIsOpen(txCreationTimestamp) {
		return nil, fmt.Errorf("ticket sell for event %s is not opened", event.Name)
	}

	sectionCapacity, err := event.GetSectionCapacity(section)
	if err != nil {
		return nil, err
	}
	currentSectionTickets, err := ticketList.CountTicketsInSection(eventID, section)
	if err != nil {
		return nil, err
	}
	if currentSectionTickets == sectionCapacity {
		return nil, fmt.Errorf("section %s is complete for event %s", section, event.Name)
	}

	ticket := NewTicket(ticketID, eventID, section, owner)
	if err = ticketList.AddTicket(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
