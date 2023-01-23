package tickenevent

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"ticken-event-contract/ledgerapi"
	"ticken-event-contract/models"
)

type eventList struct {
	stateList ledgerapi.StateList
}

func NewEventList(stub shim.ChaincodeStubInterface) EventListInterface {
	deserializeFunction := func(bytes []byte, state ledgerapi.State) error {
		return models.EventDeserialize(bytes, state.(*models.Event))
	}

	list := new(eventList)
	list.stateList = ledgerapi.NewStateList(stub, "ticken.ticket.list", deserializeFunction)

	return list
}

func (eventList *eventList) AddEvent(event *models.Event) error {
	exists, err := eventList.EventExists(event)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("event with id %s already exists", event.EventID)
	}

	return eventList.stateList.AddState(event)
}

func (eventList *eventList) UpdateEvent(event *models.Event) error {
	return eventList.stateList.UpdateState(event)
}

func (eventList *eventList) EventExists(event *models.Event) (bool, error) {
	return eventList.stateList.StateExists(event.EventID)
}

func (eventList *eventList) GetEvent(eventID string) (*models.Event, error) {
	event := new(models.Event)
	eventKey := models.EventCreateKey(eventID)

	err := eventList.stateList.GetState(eventKey, event)

	if err != nil {
		return nil, err
	}

	return event, nil
}
