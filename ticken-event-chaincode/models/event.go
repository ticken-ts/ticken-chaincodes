package models

import (
	"encoding/json"
	"fmt"
	"ticken-event-contract/ledgerapi"
	"time"
)

// *********************** Event ************************ //

type Section struct {
	Name             string `json:"name"`
	TotalTickets     int    `json:"total_tickets"`
	RemainingTickets int    `json:"remaining_tickets"`
}

type Event struct {
	EventID        string    `json:"event_id"`
	Name           string    `json:"name"`
	Date           time.Time `json:"Date"`
	Sections       []Section `json:"Sections"`
	OrganizationID string    `json:"organization_id"`
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

func NewEvent(eventID string, name string, date time.Time, organizationID string) *Event {
	event := new(Event)

	event.Name = name
	event.Date = date
	event.EventID = eventID
	event.Sections = []Section{}
	event.OrganizationID = organizationID

	return event
}

func (event *Event) getSection(name string) (*Section, bool) {
	for _, section := range event.Sections {
		if section.Name == name {
			return &section, true
		}
	}
	return nil, false
}

func (event *Event) AddSection(name string, totalTickets int) (*Section, error) {
	if totalTickets <= 0 {
		return nil, fmt.Errorf("total tickets must be greater than 0")
	}

	if event.HasSection(name) {
		return nil, fmt.Errorf("section with name %s already exists", name)
	}

	newSection := Section{
		Name:             name,
		TotalTickets:     totalTickets,
		RemainingTickets: totalTickets,
	}

	event.Sections = append(event.Sections, newSection)
	return &newSection, nil
}

func (event *Event) HasSection(name string) bool {
	_, ok := event.getSection(name)
	return ok
}

func (event *Event) IsAvailable(sectionName string) (bool, error) {
	section, ok := event.getSection(sectionName)

	if !ok {
		return false, fmt.Errorf("section does not exist")
	} else {
		return section.RemainingTickets > 0, nil
	}
}

func (event *Event) AddTicket(section string) error {
	savedSection, ok := event.getSection(section)

	if !ok {
		return fmt.Errorf("section does not exist")
	}

	if savedSection.RemainingTickets > 0 {
		savedSection.RemainingTickets -= 1
	} else {
		return fmt.Errorf("no tickets left in that section")
	}

	return nil
}

// The following implementations are required
// to implement the State.

// The following implementations are required
// to implement the State.

func (event *Event) GetKey() string {
	return ledgerapi.MakeKey(event.EventID)
}

func (event *Event) Serialize() ([]byte, error) {
	return json.Marshal(event)
}

// ****************************************************** //
