package ticken_event

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContext an interface to
// describe the minimum required functions for
// a transaction context in the commercial
// paper
type TransactionContext interface {
	contractapi.TransactionContextInterface
	GetEventList() EventListInterface
}

type SettableTransactionContext interface {
	contractapi.SettableTransactionContextInterface
}

// transactionContext implementation of
// TransactionContext for use with
// commercial paper contract
type transactionContext struct {
	contractapi.TransactionContext
	eventList EventListInterface
}

func NewTransactionContext() SettableTransactionContext {
	return new(transactionContext)
}

// GetEventList return ticken-event List
func (ctx *transactionContext) GetEventList() EventListInterface {
	if ctx.eventList == nil {
		ctx.eventList = NewEventList(ctx)
	}
	return ctx.eventList
}
