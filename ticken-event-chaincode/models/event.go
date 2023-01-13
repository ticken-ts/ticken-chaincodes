package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"ticken-event-contract/ledgerapi"
	"time"
)

// *********************** Event ************************ //

type Event struct {
	EventID  string
	Name     string
	Date     time.Time
	Sections []*Section

	// identity of the event and auditory
	MSPID             string
	OrganizerUsername string
}

func EventCreateKey(eventID string) string {
	return ledgerapi.MakeKey(eventID)
}

func EventDeserialize(jsonBytes []byte, event *Event) error {
	err := json.Unmarshal(jsonBytes, event)

	if err != nil {
		return fmt.Errorf("error deserializing ticken-event. %s", err.Error())
	}

	return nil
}

// ****************************************************** //

// ******************** Primitives ********************** //

func NewEvent(eventID uuid.UUID, name string, date time.Time, mspID string, organizerUsername string) *Event {
	return &Event{
		EventID:  eventID.String(),
		Name:     name,
		Date:     date,
		Sections: make([]*Section, 0),

		// this values will be validated from
		// the values that the chaincode notify us
		MSPID:             mspID,
		OrganizerUsername: organizerUsername,
	}
}

func (event *Event) AddSection(name string, totalTickets int, ticketPrice float64) (*Section, error) {
	if totalTickets <= 0 {
		return nil, fmt.Errorf("total tickets must be greater than 0")
	}

	if event.HasSection(name) {
		return nil, fmt.Errorf("section with name %s already exists", name)
	}

	// round the price up to two decimal
	// ex: 12.3456 -> 12.34
	twoDecimalsPrice := math.Round(ticketPrice*100) / 100

	newSection := &Section{
		Name:         name,
		SoldTickets:  0,
		TotalTickets: totalTickets,
		TicketPrice:  twoDecimalsPrice,
	}

	event.Sections = append(event.Sections, newSection)
	return newSection, nil
}

func (event *Event) HasSection(name string) bool {
	return event.findSection(name) != nil
}

func (event *Event) SectionIsAvailable(sectionName string) bool {
	section := event.findSection(sectionName)
	return section != nil && section.IsAvailable()
}

func (event *Event) SellTicketInSection(sectionName string) error {
	section := event.findSection(sectionName)

	if section == nil {
		return fmt.Errorf("section does not exist")
	}

	return section.SellTicket()
}

func (event *Event) findSection(name string) *Section {
	for _, section := range event.Sections {
		if section.Name == name {
			return section
		}
	}
	return nil
}

// The following implementations are required
// to implement the State.

func (event *Event) GetKey() string {
	return ledgerapi.MakeKey(event.EventID)
}

func (event *Event) Serialize() ([]byte, error) {
	return json.Marshal(event)
}

// ****************************************************** //
