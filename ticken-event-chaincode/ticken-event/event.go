package ticken_event

import (
	"encoding/json"
	"fmt"
	ledgerapi "ticken-ticket-contract/ledger-api"
	"time"
)

// *********************** Event ************************ //

type Section struct {
	Name             string `json:"name"`
	TotalTickets     int    `json:"total_tickets"`
	RemainingTickets int    `json:"remaining_tickets"`
}

type Event struct {
	EventID  string    `json:"event_id"`
	Name     string    `json:"name"`
	date     time.Time `json:"date"`
	sections []Section `json:"sections"`
}

func CreateEventKey(eventID string) string {
	return ledgerapi.MakeKey(eventID)
}

func Deserialize(jsonBytes []byte, event *Event) error {
	err := json.Unmarshal(jsonBytes, event)

	if err != nil {
		return fmt.Errorf("error deserializing ticken-event. %s", err.Error())
	}

	return nil
}

// ****************************************************** //

// ******************** Primitives ********************** //

func NewEvent(eventID string, name string, date time.Time) *Event {
	event := new(Event)

	event.Name = name
	event.EventID = eventID
	event.date = date
	event.sections = []Section{}

	return event
}

func (event *Event) getSection(name string) (*Section, bool) {
	for _, section := range event.sections {
		if section.Name == name {
			return &section, true
		}
	}
	return nil, false
}

func (event *Event) addSection(newSection Section) {
	event.sections = append(event.sections, newSection)
}

func (event *Event) AddSection(name string, totalTickets int) (*Section, error) {
	if totalTickets <= 0 {
		return nil, fmt.Errorf("total tickets must be greater than 0")
	}

	if event.HasSection(name) {
		return nil, fmt.Errorf("section with name %s already exists", name)
	}

	newSection := Section{Name: name, TotalTickets: totalTickets, RemainingTickets: totalTickets}
	event.addSection(newSection)
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
// to implement the StateInterface.

func (event *Event) GetSplitKey() []string {
	return []string{event.EventID}
}

func (event *Event) Serialize() ([]byte, error) {
	return json.Marshal(event)
}

// ****************************************************** //
