package ticken_event

import (
	"container/list"
	"encoding/json"
	"fmt"
	ledgerapi "ticken-ticket-contract/ledger-api"
	"time"
)

// *********************** Event ************************ //

type Section struct {
	Name         string `json:"name"`
	TotalTickets int    `json:"total_tickets"`
}

type Event struct {
	EventID  string    `json:"event_id"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
	Sections list.List `json:"sections"`
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
	event.Date = date
	event.Sections.Init()

	return event
}

func (event *Event) AddSection(name string, totalTickets int) (*Section, error) {
	if totalTickets <= 0 {
		return nil, fmt.Errorf("total tickets must be greater than 0")
	}

	if event.HasSection(name) {
		return nil, fmt.Errorf("section with name %s already exists", name)
	}

	newSection := Section{Name: name, TotalTickets: totalTickets}
	event.Sections.PushBack(newSection)
	return &newSection, nil
}

func (event *Event) HasSection(name string) bool {
	for s := event.Sections.Front(); s != nil; s = s.Next() {
		section := s.Value.(*Section)
		if section.Name == name {
			return true
		}
	}
	return false
}

// The following implementations are required
// to implement the StateInterface.

func (event *Event) GetSplitKey() []string {
	return []string{event.EventID}
}

func (event *Event) Serialize() ([]byte, error) {
	return json.Marshal(Event{})
}

// ****************************************************** //
