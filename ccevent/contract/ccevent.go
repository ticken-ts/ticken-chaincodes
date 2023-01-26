package contract

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/ticken-ts/ticken-chaincodes/common"
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

// Create a new event without any sections in the blockchain and returns
// its value. The event is created as out of sale, so it can be updated
// and sections can be added on following transactions
//
// Params
// * - eventID -> uuid format
// * - name    -> string
// * - date    -> RFC3339 format (2006-01-02T15:04:05Z07:00)
//
// The return value can be:
// * - the event created serialized in JSON format
// * - error in case some conditions to create the event are not fulfilled
func (c *Contract) Create(ctx common.ITickenTxContext, eventID, name, date string) (*Event, error) {
	_, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err // this error is already formatted
	}

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, ccErr("error parsing date: %v", err)
	}
	eventIDParsed, err := uuid.Parse(eventID)
	if err != nil {
		return nil, ccErr("error parsing event id: %v", err)
	}

	mspID, orgUsername, err := ctx.GetContextIdentity()
	if err != nil {
		return nil, ccErr("could not get context identity: %v", err)
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
		return nil, ccErr("failed to serialize event; %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return nil, ccErr("failed to updated  tate: %v", err)
	}

	return &event, nil
}

// SetEventOnSale sets the previously created event to be "on sale".
// From this moment, we can start issuing tickets for this event
//
// Params
// * - eventID -> uuid format
//
// The return value can be:
// * - error in case some conditions to create the event are not fulfilled
func (c *Contract) SetEventOnSale(ctx common.ITickenTxContext, eventID string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err // this error is already formatted
	}

	if event.OnSale {
		return ccErr("event %s already is on sale", event.EventID)
	}

	event.OnSale = true

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccErr("failed to serialize event: %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to updated  tate: %v", err)
	}

	return nil
}

// AddSection add a section on the previously created event
// Each section is identified uniquely inside the event for its name,
// so it must be unique for each section. Providing a name that is already
// in use in the same event will cause this transaction to fail
//
// Params
// * - eventID -> uuid format
// * - name    -> section name (must be unique)
// * - totalTickets
// * - ticketPrice
//
// The return value can be:
// * - the section added serialized in JSON format
// * - error in case some conditions to add the section are not fulfilled
func (c *Contract) AddSection(ctx common.ITickenTxContext, eventID, name, totalTickets, ticketPrice string) (*Section, error) {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err // this error is already formatted
	}

	totalTicketsParsed, err := strconv.Atoi(totalTickets)
	if err != nil {
		return nil, ccErr("error converting total ticket: %v", err)
	}
	ticketPriceParsed, err := strconv.ParseFloat(ticketPrice, 64)
	if err != nil {
		return nil, ccErr("error converting ticket price: %v", err)
	}

	if totalTicketsParsed <= 0 {
		return nil, ccErr("invalid total tickets value %d - total tickets must be greater than 0", totalTicketsParsed)
	}

	for _, section := range event.Sections {
		if section.Name == name {
			return nil, ccErr("section with name %s already exists", name)
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
		return nil, ccErr("failed to serialize event; %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return nil, ccErr("failed to updated the state: %v", err)
	}

	return &newSection, nil
}

func (c *Contract) GetEvent(ctx common.ITickenTxContext, eventID string) (*Event, error) {
	eventJSON, err := ctx.GetStub().GetState(eventID)
	if err != nil {
		return nil, ccErr("failed to read event: %v", err)
	}
	if eventJSON == nil {
		return nil, ccErr("the event %s does not exist", eventID)
	}

	var event Event
	if err := json.Unmarshal(eventJSON, &event); err != nil {
		return nil, ccErr("failed to deserialize event: %v", err)
	}

	return &event, nil
}

func (c *Contract) SellTicket(ctx common.ITickenTxContext, eventID string, sectionName string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err // this error is already formatted
	}

	if !event.OnSale {
		return ccErr("event not on sale")
	}

	var foundSection *Section
	for _, section := range event.Sections {
		if section.Name == sectionName {
			foundSection = section
			break
		}
	}

	if foundSection == nil {
		return ccErr("section %s doest not exist in event %s", sectionName, eventID)
	}

	if foundSection.SoldTickets == foundSection.TotalTickets {
		return ccErr("section %s is full", sectionName)
	}

	foundSection.SoldTickets += 1
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccErr("failed to deserialize event: %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to update ledger: %v", err)
	}

	return nil
}

func ccErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args)
	return fmt.Errorf("[%s] | %v", Name, msg)
}
