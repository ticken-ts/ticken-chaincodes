package ticken_ticket

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the commercial
// paper
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetTicketList() ListInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// commercial paper contract
type TransactionContext struct {
	contractapi.TransactionContext
	ticketList *List
}

// GetTicketList return ticken-event List
func (ctx *TransactionContext) GetTicketList() ListInterface {
	if ctx.ticketList == nil {
		ctx.ticketList = newList(ctx)
	}

	return ctx.ticketList
}
