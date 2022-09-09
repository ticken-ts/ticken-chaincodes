package tickenticket

import "ticken-ticket-contract/ledgerapi"

type ticketList struct {
	stateList ledgerapi.StateListInterface
}

func NewTicketList(ctx TransactionContextInterface) *ticketList {
	stateList := new(ledgerapi.StateList)

	stateList.Ctx = ctx
	stateList.Name = "ticken.ticken-event.ticketList"

	stateList.Deserialize = func(bytes []byte, state ledgerapi.State) error {
		return TicketDeserialize(bytes, state.(*Ticket))
	}

	list := new(ticketList)
	list.stateList = stateList

	return list
}

func (ticketList *ticketList) AddTicket(ticket *Ticket) error {
	return ticketList.stateList.AddState(ticket)
}

func (ticketList *ticketList) UpdateTicket(ticket *Ticket) error {
	return ticketList.stateList.UpdateState(ticket)
}

func (ticketList *ticketList) GetTicket(eventID string, ticketID string) (*Ticket, error) {
	ticket := new(Ticket)
	ticketKey := TicketCreateKey(eventID, ticketID)

	if err := ticketList.stateList.GetState(ticketKey, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
