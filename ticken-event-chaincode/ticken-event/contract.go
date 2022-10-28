package ticken_event

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"ticken-event-contract/models"
	"time"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Create(ctx ITickenTxContext, eventID string, name string, date string) (*models.Event, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, fmt.Errorf("error parsing Date %s", err.Error())
	}

	organizationID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("could not get organization id")
	}

	newEvent := models.NewEvent(eventID, name, parsedDate, organizationID)
	err = ctx.GetEventList().AddEvent(newEvent)
	if err != nil {
		return nil, err
	}

	err = ctx.GetNotifier().NotifyEventCreation(newEvent)
	if err != nil {
		return nil, err
	}
	return newEvent, nil
}

func (c *Contract) AddSection(ctx ITickenTxContext, eventID string, name string, totalTickets string) (*models.Event, error) {
	totalTicketsParsed, conversionError := strconv.Atoi(totalTickets)
	if conversionError != nil {
		return nil, conversionError
	}

	event, getEventError := ctx.GetEventList().GetEvent(eventID)
	if getEventError != nil {
		return nil, getEventError
	}

	_, addSectionError := event.AddSection(name, totalTicketsParsed)
	if addSectionError != nil {
		return nil, addSectionError
	}

	if updateError := ctx.GetEventList().UpdateEvent(event); updateError != nil {
		return nil, updateError
	}

	err := ctx.GetNotifier().NotifySectionAddition(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Contract) Get(ctx ITickenTxContext, eventID string) (*models.Event, error) {
	event, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Contract) EventExists(ctx ITickenTxContext, eventID string) (bool, error) {
	_, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return false, nil
	} else {
		return true, nil
	}
}

func (c *Contract) IsAvailable(ctx ITickenTxContext, eventID string, section string) (bool, error) {
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

func (c *Contract) AddTicket(ctx ITickenTxContext, eventID string, section string) error {
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
