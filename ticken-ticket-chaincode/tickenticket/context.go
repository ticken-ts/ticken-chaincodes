package tickenticket

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/ccinvokers"
	"ticken-ticket-contract/ccinvokers/tickenevent"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the commercial
// paper
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetTicketList() ListInterface
	GetTickenEventCCInvoker() ccinvokers.TickenEventInvoker
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// commercial paper contract
type transactionContext struct {
	contractapi.TransactionContext
	ticketList           ListInterface
	tickenEventCCInvoker ccinvokers.TickenEventInvoker
}

func NewTransactionContext() *transactionContext {
	return new(transactionContext)
}

// GetTicketList return ticken-event ticketList
func (ctx *transactionContext) GetTicketList() ListInterface {
	if ctx.ticketList == nil {
		ctx.ticketList = NewTicketList(ctx)
	}
	return ctx.ticketList
}

func (ctx *transactionContext) GetTickenEventCCInvoker() ccinvokers.TickenEventInvoker {
	if ctx.tickenEventCCInvoker == nil {
		ctx.tickenEventCCInvoker = tickenevent.NewInvoker(ctx.GetStub())
	}
	return ctx.tickenEventCCInvoker
}
