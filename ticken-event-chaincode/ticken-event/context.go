package ticken_event

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the commercial
// paper
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetEventList() EventListInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// commercial paper contract
type TransactionContext struct {
	contractapi.TransactionContext
	eventList *eventList
}

// GetEventList return ticken-event List
func (ctx *TransactionContext) GetEventList() EventListInterface {
	if ctx.eventList == nil {
		ctx.eventList = NewEventList(ctx)
	}
	return ctx.eventList
}
