package ticken_ticket

import (
	ledgerapi "ticken-ticket-contract/ledger-api"
)

// *********************** List ************************* //

type ListInterface interface {
	AddTicket(ticket *Ticket) error
	UpdateTicket(ticket *Ticket) error
	GetTicket(eventID string, ticketID string) (*Ticket, error)
}

type List struct {
	stateList ledgerapi.StateListInterface
}

func newList(ctx TransactionContextInterface) *List {
	stateList := new(ledgerapi.StateList)

	stateList.Ctx = ctx
	stateList.Name = "ticken.ticken-event.list"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*Ticket))
	}

	list := new(List)
	list.stateList = stateList

	return list
}

// ****************************************************** //

// ******************** Primitives ********************** //

func (ticketList *List) AddTicket(ticket *Ticket) error {
	return ticketList.stateList.AddState(ticket)
}

func (ticketList *List) UpdateTicket(ticket *Ticket) error {
	return ticketList.stateList.UpdateState(ticket)
}

func (ticketList *List) GetTicket(eventID string, ticketID string) (*Ticket, error) {
	ticket := new(Ticket)
	ticketKey := CreateTicketKey(eventID, ticketID)

	err := ticketList.stateList.GetState(ticketKey, ticket)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

// ****************************************************** //
