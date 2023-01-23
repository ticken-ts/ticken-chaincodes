package contract

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"math"
	"strconv"
	"time"
)

type Contract struct {
	contractapi.Contract
}

const Name = "cc-event"

type Event struct {
	EventID  uuid.UUID  `json:"event_id"`
	Name     string     `json:"name"`
	Date     time.Time  `json:"date"`
	Sections []*Section `json:"sections"`

	// indicates if the event is currently
	// available to sell tickets
	OnSale bool `json:"on_sale"`

	// identity of the event and auditory
	MSPID             string `json:"msp_id"`
	OrganizerUsername string `json:"organizer_username"`
}

type Section struct {
	EventID      uuid.UUID `json:"event_id"`
	Name         string    `json:"name"`
	TicketPrice  float64   `json:"ticket_price"`
	TotalTickets int       `json:"total_tickets"`
	SoldTickets  int       `json:"sold_tickets"`
}

func (c *Contract) Create(ctx contractapi.TransactionContextInterface, eventID, name, date string) (*Event, error) {
	_, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, ccError(fmt.Errorf("error parsing date: %v", err))
	}

	eventIDParsed, err := uuid.Parse(eventID)
	if err != nil {
		return nil, ccError(fmt.Errorf("error parsing event id: %v", err))
	}

	mspID, orgUsername, err := getContextIdentity(ctx)
	if err != nil {
		return nil, ccError(fmt.Errorf("cound not get context identity: %v", err))
	}

	event := Event{
		EventID:  eventIDParsed,
		Name:     name,
		Date:     parsedDate,
		Sections: make([]*Section, 0),

		// initially the event is marked to
		// be out of sale, so it can be updated
		OnSale: false,

		// this values will be validated from
		// the values that the chaincode notify us
		MSPID:             mspID,
		OrganizerUsername: orgUsername,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, ccError(fmt.Errorf("failed to serialize event; %v", err))
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return nil, ccError(fmt.Errorf("failed to updated the state: %v", err))
	}

	return &event, nil
}

func (c *Contract) AddSection(ctx contractapi.TransactionContextInterface, eventID, name, totalTickets, ticketPrice string) (*Section, error) {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	totalTicketsParsed, err := strconv.Atoi(totalTickets)
	if err != nil {
		return nil, ccError(fmt.Errorf("error converting total ticket: %v", err))
	}
	ticketPriceParsed, err := strconv.ParseFloat(ticketPrice, 64)
	if err != nil {
		return nil, ccError(fmt.Errorf("error converting ticket price: %v", err))
	}

	if totalTicketsParsed <= 0 {
		return nil, ccError(fmt.Errorf("invalid total tickets value %d - total tickets must be greater than 0", totalTicketsParsed))
	}

	for _, section := range event.Sections {
		if section.Name == name {
			return nil, ccError(fmt.Errorf("section with name %s already exists", name))
		}
	}

	// round the price up to two decimal
	// ex: 12.3456 -> 12.34
	twoDecimalsPrice := math.Round(ticketPriceParsed*100) / 100

	newSection := Section{
		Name:         name,
		EventID:      event.EventID,
		SoldTickets:  0,
		TotalTickets: totalTicketsParsed,
		TicketPrice:  twoDecimalsPrice,
	}

	event.Sections = append(event.Sections, &newSection)

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize event; %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return nil, ccError(fmt.Errorf("failed to updated the state: %v", err))
	}

	return &newSection, nil
}

func (c *Contract) GetEvent(ctx contractapi.TransactionContextInterface, eventID string) (*Event, error) {
	eventJSON, err := ctx.GetStub().GetState(eventID)
	if err != nil {
		return nil, ccError(fmt.Errorf("failed to read event: %v", err))
	}
	if eventJSON == nil {
		return nil, ccError(fmt.Errorf("the event %s does not exist", eventID))
	}

	var event Event
	err = json.Unmarshal(eventJSON, &event)
	if err != nil {
		return nil, ccError(fmt.Errorf("failed to descerialize event: %v", err))
	}

	return &event, nil
}

func (c *Contract) SellTicket(ctx contractapi.TransactionContextInterface, eventID string, sectionName string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}

	if !event.OnSale {
		return ccError(fmt.Errorf("event not on sale"))
	}

	var foundSection *Section
	for _, section := range event.Sections {
		if section.Name == sectionName {
			foundSection = section
			break
		}
	}

	if foundSection == nil {
		return ccError(fmt.Errorf("section %s doest not exist in event %s", sectionName, eventID))
	}

	if foundSection.SoldTickets == foundSection.TotalTickets {
		return ccError(fmt.Errorf("section %s is full", sectionName))
	}

	foundSection.SoldTickets += 1
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccError(fmt.Errorf("failed to descerialize event: %v", err))
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccError(fmt.Errorf("failed to upadate ledger: %v", err))
	}

	return nil
}

func getContextIdentity(ctx contractapi.TransactionContextInterface) (string, string, error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", "", err
	}

	username, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", "", err
	}

	return mspID, username, nil
}

func ccError(err error) error {
	return fmt.Errorf("[%s] | %v", Name, err)
}
