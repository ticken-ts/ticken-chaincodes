package tickenticket

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/ccinvokers"
)

// TransactionContext implementation of
// TickenTxContext for use with
// commercial paper contract
type tickenTxContext struct {
	contractapi.TransactionContext
	ticketList           TicketList
	tickenEventCCInvoker ccinvokers.TickenEventInvoker
}

func NewTransactionContext() *tickenTxContext {
	return new(tickenTxContext)
}

// GetTicketList return ticken-event ticketList
func (ctx *tickenTxContext) GetTicketList() TicketList {
	if ctx.ticketList == nil {
		ctx.ticketList = NewTicketList(ctx.GetStub())
	}
	return ctx.ticketList
}

func (ctx *tickenTxContext) GetTickenEventInvoker() ccinvokers.TickenEventInvoker {
	if ctx.tickenEventCCInvoker == nil {
		ctx.tickenEventCCInvoker = ccinvokers.NewTickenEventInvoker(ctx.GetStub())
	}
	return ctx.tickenEventCCInvoker
}
