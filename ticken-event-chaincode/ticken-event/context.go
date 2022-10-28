package ticken_event

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-event-contract/ccnotifier"
)

// TransactionContext an interface to
// describe the minimum required functions for
// a transaction context in the commercial
// paper

type TickenTxContext struct {
	contractapi.TransactionContext
	eventList EventListInterface
	notifier  ccnotifier.Notifier
}

func NewTransactionContext() *TickenTxContext {
	return new(TickenTxContext)
}

// GetEventList return ticken-event List
func (ctx *TickenTxContext) GetEventList() EventListInterface {
	if ctx.eventList == nil {
		ctx.eventList = NewEventList(ctx.GetStub())
	}
	return ctx.eventList
}

// GetEventList return cc-notifier
func (ctx *TickenTxContext) GetNotifier() ccnotifier.Notifier {
	if ctx.eventList == nil {
		ctx.notifier = ccnotifier.NewNotifier(ctx.GetStub())
	}
	return ctx.notifier
}
