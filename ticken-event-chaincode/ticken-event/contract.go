package ticken_event

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) InitSamples(ctx TransactionContext) error {
	events := [3]*Event{}
	events[0] = NewEvent("62dc7486b6faaccaf51089d9", "Event 1", time.Now().AddDate(0, 0, 20))
	events[1] = NewEvent("62dc75325721c6ec1dda26e6", "Event 2", time.Now().AddDate(0, 0, 20))
	events[2] = NewEvent("62dc753a7ab2e97c2afb0f6b", "Event 3", time.Now().AddDate(0, 0, 20))

	eventList := ctx.GetEventList()
	for i := 0; i < len(events); i++ {
		err := eventList.AddEvent(events[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Contract) Get(ctx TransactionContext, eventID string) (*Event, error) {
	event, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Contract) Create(ctx TransactionContext, eventID string, name string, date string) (*Event, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date %s", err.Error())
	}

	newEvent := NewEvent(eventID, name, parsedDate)

	err = ctx.GetEventList().AddEvent(newEvent)
	if err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (c *Contract) AddSection(ctx TransactionContext, eventID string, name string, totalTickets int) (*Event, error) {
	event, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	_, err = event.AddSection(name, totalTickets)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Contract) EventExists(ctx TransactionContext, eventID string) (bool, error) {
	return true, nil
}

func (c *Contract) IsAvailable(ctx TransactionContext, eventID string, section string) (bool, error) {
	return true, nil
}

func (c *Contract) AddTicket(ctx TransactionContext, eventID string, section string) error {
	return nil
}
