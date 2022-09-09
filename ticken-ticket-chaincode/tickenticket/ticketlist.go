package tickenticket

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"ticken-ticket-contract/ledgerapi"
)

type ticketList struct {
	stateList ledgerapi.StateListInterface
}

func NewTicketList(stub shim.ChaincodeStubInterface) *ticketList {
	deserializeFunction := func(bytes []byte, state ledgerapi.State) error {
		return TicketDeserialize(bytes, state.(*Ticket))
	}

	list := new(ticketList)
	list.stateList = ledgerapi.NewStateList(stub, "ticken.ticket.list", deserializeFunction)

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
