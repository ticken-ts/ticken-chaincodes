package tickenticket

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/ccinvokers"
)

// TransactionContext implementation of
// ITickenTxContext for use with
// commercial paper contract
type TickenTxContext struct {
	contractapi.TransactionContext
	ticketList           TicketList
	tickenEventCCInvoker ccinvokers.TickenEventInvoker
}

func NewTransactionContext() *TickenTxContext {
	return new(TickenTxContext)
}

// GetTicketList return ticken-event ticketList
func (ctx *TickenTxContext) GetTicketList() TicketList {
	if ctx.ticketList == nil {
		ctx.ticketList = NewTicketList(ctx.GetStub())
	}
	return ctx.ticketList
}

func (ctx *TickenTxContext) GetTickenEventInvoker() ccinvokers.TickenEventInvoker {
	if ctx.tickenEventCCInvoker == nil {
		ctx.tickenEventCCInvoker = ccinvokers.NewTickenEventInvoker(ctx.GetStub())
	}
	return ctx.tickenEventCCInvoker
}
