package contract

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
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

func (c *Contract) Create(ctx common.ITickenTxContext, eventID, name, date string) peer.Response {
	_, err := findEvent(ctx, eventID)
	if err != nil {
		return ccErr(err.Error())
	}

	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ccErr("error parsing date: %v", err)
	}
	eventIDParsed, err := uuid.Parse(eventID)
	if err != nil {
		return ccErr("error parsing event id: %v", err)
	}

	mspID, orgUsername, err := ctx.GetContextIdentity()
	if err != nil {
		return ccErr("could not get context identity: %v", err)
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
		return ccErr("failed to serialize event; %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to updated  tate: %v", err)
	}

	return shim.Success(eventJSON)
}

func (c *Contract) SetEventOnSale(ctx common.ITickenTxContext, eventID string) peer.Response {
	event, err := findEvent(ctx, eventID)
	if err != nil {
		return ccErr(err.Error())
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

	return shim.Success(nil)
}

func (c *Contract) AddSection(ctx common.ITickenTxContext, eventID, name, totalTickets, ticketPrice string) peer.Response {
	event, err := findEvent(ctx, eventID)
	if err != nil {
		return ccErr(err.Error())
	}

	totalTicketsParsed, err := strconv.Atoi(totalTickets)
	if err != nil {
		return ccErr("error converting total ticket: %v", err)
	}
	ticketPriceParsed, err := strconv.ParseFloat(ticketPrice, 64)
	if err != nil {
		return ccErr("error converting ticket price: %v", err)
	}

	if totalTicketsParsed <= 0 {
		return ccErr("invalid total tickets value %d - total tickets must be greater than 0", totalTicketsParsed)
	}

	for _, section := range event.Sections {
		if section.Name == name {
			return ccErr("section with name %s already exists", name)
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
		return ccErr("failed to serialize event; %v", err)
	}

	if err := ctx.GetStub().PutState(eventID, eventJSON); err != nil {
		return ccErr("failed to updated the state: %v", err)
	}

	sectionJSON, err := json.Marshal(&newSection)
	if err != nil {
		return ccErr("failed to serialize section; %v", err)
	}

	return shim.Success(sectionJSON)
}

func (c *Contract) GetEvent(ctx common.ITickenTxContext, eventID string) peer.Response {
	event, err := findEvent(ctx, eventID)
	if err != nil {
		return ccErr(err.Error())
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return ccErr("failed to serialize event: %v", err)
	}

	return shim.Success(eventJSON)
}

func (c *Contract) SellTicket(ctx common.ITickenTxContext, eventID string, sectionName string) peer.Response {
	event, err := findEvent(ctx, eventID)
	if err != nil {
		return ccErr(err.Error())
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

	return shim.Success(nil)
}

func findEvent(ctx common.ITickenTxContext, eventID string) (*Event, error) {
	eventJSON, err := ctx.GetStub().GetState(eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to read event: %v", err)
	}
	if eventJSON == nil {
		return nil, fmt.Errorf("the event %s does not exist", eventID)
	}

	var event Event
	if err := json.Unmarshal(eventJSON, &event); err != nil {
		return nil, fmt.Errorf("failed to deserialize event: %v", err)
	}

	return &event, nil
}

func ccErr(format string, args ...any) peer.Response {
	msg := fmt.Sprintf(format, args)
	return shim.Error(fmt.Sprintf("[%s] | %v", Name, msg))
}
