package ticken_event

import (
	"fmt"
	ledgerapi "ticken-ticket-contract/ledger-api"
)

// *********************** List ************************* //

type EventListInterface interface {
	AddEvent(event *Event) error
	UpdateEvent(event *Event) error
	EventExists(event *Event) bool
	GetEvent(evenID string) (*Event, error)
}

type eventList struct {
	stateList ledgerapi.StateListInterface
}

// ****************************************************** //

// ******************** Primitives ********************** //

func NewEventList(ctx TransactionContextInterface) *eventList {
	stateList := new(ledgerapi.StateList)

	stateList.Ctx = ctx
	stateList.Name = "ticken.ticken-event.list"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*Event))
	}

	list := new(eventList)
	list.stateList = stateList

	return list
}

func (eventList *eventList) AddEvent(event *Event) error {
	if eventList.EventExists(event) {
		return fmt.Errorf("event with id %s already exists", event.EventID)
	}

	return eventList.stateList.AddState(event)
}

func (eventList *eventList) UpdateEvent(event *Event) error {
	return eventList.stateList.UpdateState(event)
}

func (eventList *eventList) EventExists(event *Event) bool {
	return eventList.stateList.KeyExits(event.EventID)
}

func (eventList *eventList) GetEvent(eventID string) (*Event, error) {
	event := new(Event)
	eventKey := CreateEventKey(eventID)

	err := eventList.stateList.GetState(eventKey, event)

	if err != nil {
		return nil, err
	}

	return event, nil
}

// ****************************************************** //
