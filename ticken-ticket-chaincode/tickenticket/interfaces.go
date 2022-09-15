package tickenticket

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/ccinvokers"
)

type TicketList interface {
	AddTicket(ticket *Ticket) error
	UpdateTicket(ticket *Ticket) error
	TicketExist(eventID string, ticketID string) (bool, error)
	GetTicket(eventID string, ticketID string) (*Ticket, error)
}

type TickenTxContext interface {
	contractapi.TransactionContextInterface
	GetTicketList() TicketList
	GetTickenEventInvoker() ccinvokers.TickenEventInvoker
}
