package tickenevent

import (
	"fmt"
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

// GetNotifier return cc-notifier
func (ctx *TickenTxContext) GetNotifier() ccnotifier.Notifier {
	if ctx.notifier == nil {
		ctx.notifier = ccnotifier.NewNotifier(ctx.GetStub())
	}
	return ctx.notifier
}

func (ctx *TickenTxContext) GetCallerIdentity() (string, string, error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", "", fmt.Errorf("could not get MSP ID")
	}

	username, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", "", fmt.Errorf("could not get user")
	}

	return mspID, username, nil
}
