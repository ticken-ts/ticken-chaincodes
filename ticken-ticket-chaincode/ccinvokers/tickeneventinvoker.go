package ccinvokers

import (
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

const (
	TickenEventName = "ticken-event"

	// functions
	GetFunction = "Get"
)

type tickenEventInvoker struct {
	baseInvoker BaseInvoker
}

func NewTickenEventInvoker(stub shim.ChaincodeStubInterface) *tickenEventInvoker {
	return &tickenEventInvoker{
		baseInvoker: NewBaseInvoker(TickenEventName, stub),
	}
}

func (i *tickenEventInvoker) GetEvent(eventID string) (*Event, error) {
	payload, err := i.baseInvoker.Invoke(GetFunction, eventID)
	if err != nil {
		return nil, err
	}

	var event = new(Event)
	err = json.Unmarshal(payload, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
