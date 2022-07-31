package ticken_ticket

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	ticken_event_cc "ticken-ticket-contract/cc-invokers/ticken-event-cc"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the commercial
// paper
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetTicketList() ListInterface
	GetTickenEventCCInvoker() ticken_event_cc.TickenEventCCInvoker
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// commercial paper contract
type TransactionContext struct {
	contractapi.TransactionContext
	ticketList           ListInterface
	tickenEventCCInvoker ticken_event_cc.TickenEventCCInvoker
}

// GetTicketList return ticken-event List
func (ctx *TransactionContext) GetTicketList() ListInterface {
	if ctx.ticketList == nil {
		ctx.ticketList = newTicketList(ctx)
	}

	return ctx.ticketList
}

func (ctx *TransactionContext) GetTickenEventCCInvoker() ticken_event_cc.TickenEventCCInvoker {
	if ctx.tickenEventCCInvoker == nil {
		ctx.tickenEventCCInvoker = ticken_event_cc.NewTickenEventCCInvoker(ctx)
	}

	return ctx.tickenEventCCInvoker
}
