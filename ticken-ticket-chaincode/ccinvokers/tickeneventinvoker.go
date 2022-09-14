package ccinvokers

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

const (
	TickenEventName = "ticken-event"
)

type tickenEventInvoker struct {
	baseInvoker BaseInvoker
}

func NewTickenEventInvoker(stub shim.ChaincodeStubInterface) *tickenEventInvoker {
	return &tickenEventInvoker{
		baseInvoker: NewBaseInvoker(TickenEventName, stub),
	}
}

func (i *tickenEventInvoker) IsAvailable(eventID string, section string) (bool, error) {
	return true, nil
}

func (i *tickenEventInvoker) AddTicket(eventID string, section string) error {
	return nil
}
