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

type EventStatus string

const (
	// EventStatusDraft is the status of an
	// event that is not yet published
	EventStatusDraft EventStatus = "draft"

	// EventStatusOnSale is the status of an
	// event that is published for sale
	EventStatusOnSale EventStatus = "on_sale"

	// EventStatusRunning is the status of an
	// event that is currently happening
	EventStatusRunning EventStatus = "running"

	// EventStatusFinished is the status of an
	// event that has finished
	EventStatusFinished EventStatus = "finished"
)

type Event struct {
	EventID  string      `json:"event_id"`
	Name     string      `json:"name"`
	Date     time.Time   `json:"date"`
	Sections []*Section  `json:"sections"`
	Status   EventStatus `json:"status"`

	// identity of the event and auditory
	MSPID             string `json:"msp_id"`
	OrganizerUsername string `json:"organizer_username"`
}

type Section struct {
	EventID      string  `json:"event_id"`
	Name         string  `json:"name"`
	TicketPrice  float64 `json:"ticket_price"`
	TotalTickets int     `json:"total_tickets"`
	SoldTickets  int     `json:"sold_tickets"`
}

// Create a new event without any sections in the blockchain and returns
// its value. The event is created as with status "EventStatusDraft", so it
// can be updated and sections can be added on following transactions
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
	existentEvent, err := c.GetEvent(ctx, eventID)
	if existentEvent != nil {
		return nil, ccErr("event with ID %s already exists", eventID)
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
		EventID:  eventIDParsed.String(),
		Name:     name,
		Date:     parsedDate,
		Sections: make([]*Section, 0),
		Status:   EventStatusDraft,

		// this values will be validated from
		// the values that the chaincode notify us
		MSPID:             mspID,
		OrganizerUsername: orgUsername,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, ccErr("failed to serialize event: %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return nil, ccErr("failed to updated  tate: %v", err)
	}

	return &event, nil
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

	if event.Status != EventStatusDraft {
		return nil, ccErr("event is not in status draft")
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

// Sell sets the previously created event to be on status
// "on sale".  From this moment, we can start issuing tickets for
// this event. In addition, this status blocks any modification or change
// in the event, including adding sections
//
// Params
// * - eventID -> uuid format
//
// The return value can be:
// * - error in case the event cant transition to state "on_sale"
func (c *Contract) Sell(ctx common.ITickenTxContext, eventID string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err // this error is already formatted
	}

	if event.Status == EventStatusOnSale {
		return ccErr("event %s already is on status", event.EventID, EventStatusOnSale)
	}

	if event.Status != EventStatusDraft {
		return ccErr("event cant go from %s to %s", event.Status, EventStatusDraft)
	}

	// update status from
	// EventStatusOnSale -> EventStatusOnSale
	event.Status = EventStatusOnSale

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccErr("failed to serialize event: %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to updated  tate: %v", err)
	}

	return nil
}

// Start sets the previously created event to be on status
// "running".  From this moment, is not possible to issue additional
// tickets for the event and all tickets start to become available to
// be scanned
//
// Params
// * - eventID -> uuid format
//
// The return value can be:
// * - error in case the event cant transition to state "running"
func (c *Contract) Start(ctx common.ITickenTxContext, eventID string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err // this error is already formatted
	}

	if event.Status == EventStatusRunning {
		return ccErr("event %s already is on status %s", event.EventID, EventStatusRunning)
	}

	if event.Status != EventStatusRunning {
		return ccErr("event cant go from %s to %s", event.Status, EventStatusDraft)
	}

	// update status from
	// EventStatusOnSale -> EventStatusOnSale
	event.Status = EventStatusRunning

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccErr("failed to serialize event: %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to updated  tate: %v", err)
	}

	return nil
}

// Finish sets the previously created event to be on status
// "finished". Once in this state, all tickets will be invalidated
// and they will be free to trade on public blockchain as collectibles,
// or in other words, without any extra cost.
//
// Params
// * - eventID -> uuid format
//
// The return value can be:
// * - error in case the event cant transition to state "running"
func (c *Contract) Finish(ctx common.ITickenTxContext, eventID string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err // this error is already formatted
	}

	if event.Status == EventStatusFinished {
		return ccErr("event %s already is on status %s", event.EventID, EventStatusFinished)
	}

	if event.Status != EventStatusRunning {
		return ccErr("event cant go from %s to %s", event.Status, EventStatusRunning)
	}

	// update status from
	// EventStatusRunning -> EventStatusRunning
	event.Status = EventStatusRunning

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccErr("failed to serialize event: %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to updated  tate: %v", err)
	}

	return nil
}

// GetEvent returns the event information of the event with id "eventID".
//
// Params
// * - eventID -> uuid format
//
// The return value can be:
// * - error in case of the event is not found
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

// SellTicket increase in one the ticket count on the
// section with "sectionName" of the event with "eventID"
// The event must be in the state "OnSale" in order to success
//
// Params
// * - eventID -> uuid format
// * - sectionName -> unique name that identifies the section in the event
//
// The return value can be:
// * - error in case of the event is not found
func (c *Contract) SellTicket(ctx common.ITickenTxContext, eventID string, sectionName string) error {
	event, err := c.GetEvent(ctx, eventID)
	if err != nil {
		return err // this error is already formatted
	}

	if event.Status != EventStatusOnSale {
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
