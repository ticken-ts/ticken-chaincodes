package ticken_event

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"ticken-event-contract/models"
	"time"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Create(ctx ITickenTxContext, eventID, name, date string) (*models.Event, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, fmt.Errorf("error parsing Date %s", err.Error())
	}
	uuidEventID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, fmt.Errorf("error parsing event id %s", err.Error())
	}

	mspID, organizerUsername, err := ctx.GetCallerIdentity()
	if err != nil {
		return nil, err
	}

	newEvent := models.NewEvent(uuidEventID, name, parsedDate, mspID, organizerUsername)

	if err = ctx.GetEventList().AddEvent(newEvent); err != nil {
		return nil, err
	}

	if err = ctx.GetNotifier().NotifyEventCreation(newEvent); err != nil {
		return nil, err
	}

	return newEvent, nil
}

func (c *Contract) AddSection(ctx ITickenTxContext, eventID, name, totalTickets, ticketPrice string) (*models.Event, error) {
	totalTicketsParsed, err := strconv.Atoi(totalTickets)
	if err != nil {
		return nil, fmt.Errorf("error converting total ticket")
	}
	ticketPriceParsed, err := strconv.ParseFloat(ticketPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting ticket price")
	}

	event, err := ctx.GetEventList().GetEvent(eventID)
	if err != nil {
		return nil, err
	}

	newSection, err := event.AddSection(name, totalTicketsParsed, ticketPriceParsed)
	if err != nil {
		return nil, err
	}

	if err = ctx.GetEventList().UpdateEvent(event); err != nil {
		return nil, err
	}

	if err = ctx.GetNotifier().NotifySectionAddition(newSection, eventID); err != nil {
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

	isAvailable := event.SectionIsAvailable(section)

	return isAvailable, nil
}

func (c *Contract) AddTicket(ctx ITickenTxContext, eventID string, section string) error {
	event, err := ctx.GetEventList().GetEvent(eventID)

	if err != nil {
		return err
	}

	err = event.SellTicketInSection(section)

	if err != nil {
		return err
	}

	err = ctx.GetEventList().UpdateEvent(event)

	if err != nil {
		return err
	}

	return nil
}
