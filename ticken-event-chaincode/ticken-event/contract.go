package ticken_event

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
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

func (c *Contract) Create(ctx TransactionContext, eventID string, name string, date string) (*EventDTO, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, fmt.Errorf("error parsing Date %s", err.Error())
	}

	newEvent := NewEvent(eventID, name, parsedDate)

	err = ctx.GetEventList().AddEvent(newEvent)
	if err != nil {
		return nil, err
	}

	organizationID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, err
	}

	eventDTO := NewEventDTO(newEvent, organizationID)
	serialized, err := eventDTO.Serialize()
	if err != nil {
		return nil, err
	}

	if err := ctx.GetStub().SetEvent("create", serialized); err != nil {
		return nil, err
	}
	return eventDTO, nil
}

func (c *Contract) AddSection(ctx TransactionContext, eventID string, name string, totalTickets string) error {

	totalTicketsParsed, conversionError := strconv.Atoi(totalTickets)
	if conversionError != nil {
		return conversionError
	}

	event, getEventError := ctx.GetEventList().GetEvent(eventID)
	if getEventError != nil {
		return getEventError
	}

	_, addSectionError := event.AddSection(name, totalTicketsParsed)
	if addSectionError != nil {
		return addSectionError
	}

	if updateError := ctx.GetEventList().UpdateEvent(event); updateError != nil {
		return updateError
	}

	return nil
}

func (c *Contract) EventExists(ctx TransactionContext, eventID string) (bool, error) {
	_, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (c *Contract) IsAvailable(ctx TransactionContext, eventID string, section string) (bool, error) {
	event, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return false, err
	}
	isAvailable, isAvailableErr := event.IsAvailable(section)

	if isAvailableErr != nil {
		return false, isAvailableErr
	}

	return isAvailable, nil
}

func (c *Contract) AddTicket(ctx TransactionContext, eventID string, section string) error {
	event, err := ctx.GetEventList().GetEvent(eventID)

	if err != nil {
		return err
	}

	err = event.AddTicket(section)

	if err != nil {
		return err
	}

	err = ctx.GetEventList().UpdateEvent(event)

	if err != nil {
		return err
	}

	return nil
}
