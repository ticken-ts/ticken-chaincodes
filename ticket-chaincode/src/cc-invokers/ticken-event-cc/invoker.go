package ticken_event_cc

import (
	"encoding/json"
	ccinvokers "ticken-ticket-contract/cc-invokers"
	tickenticket "ticken-ticket-contract/ticken-ticket"
)

type TickenEventCCInvoker interface {
	GetEvent(eventID string) (*Event, error)
}

const (
	TickenEventCCName = "ticken-event"

	// functions
	GetFunction = "Get"
)

type tickenEventCCInvoker struct {
	ccBaseInvoker ccinvokers.CCBaseInvoker
}

func NewTickenEventCCInvoker(ctx tickenticket.TransactionContextInterface) TickenEventCCInvoker {
	tickenEventCCInvoker := new(tickenEventCCInvoker)
	tickenEventCCInvoker.ccBaseInvoker = ccinvokers.NewCCBaseInvoker(TickenEventCCName, ctx)
	return tickenEventCCInvoker
}

func (i *tickenEventCCInvoker) GetEvent(eventID string) (*Event, error) {
	payload, err := i.ccBaseInvoker.Invoke(GetFunction, eventID)
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
