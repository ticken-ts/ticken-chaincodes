package ticken_ticket

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

func (c *Contract) Issue(ctx TransactionContextInterface, ticketID string, eventID string, owner string) (*Ticket, error) {
	ticketList := ctx.GetTicketList()

	ticket := new(Ticket)
	ticket.Init(ticketID, eventID, owner)

	err := ticketList.AddTicket(ticket)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (c *Contract) Sign(ctx TransactionContextInterface, eventName string, ticketID string) (*Ticket, error) {
	ticket, err := ctx.GetTicketList().GetTicket(eventName, ticketID)

	if err != nil {
		return nil, err
	}

	ticket.Sign()
	return ticket, nil
}

func (c *Contract) Scan(ctx TransactionContextInterface, eventName string, ticketID string) (*Ticket, error) {
	ticket, err := ctx.GetTicketList().GetTicket(eventName, ticketID)

	if err != nil {
		return nil, err
	}

	ticket.Scan()
	return ticket, nil
}
