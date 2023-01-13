package ccnotifier

import (
	"ticken-event-contract/models"
	"time"
)

type SectionDTO struct {
	EventID      string `json:"event_id"`
	Name         string `json:"name"`
	TotalTickets int    `json:"total_tickets"`
	SoldTickets  int    `json:"sold_tickets"`
}

type EventDTO struct {
	EventID           string    `json:"event_id"`
	Name              string    `json:"name"`
	Date              time.Time `json:"date"`
	MSPID             string    `json:"msp_id"`
	OrganizerUsername string    `json:"organizer_username"`
}

func MapEventToDTO(event *models.Event) *EventDTO {
	return &EventDTO{
		EventID:           event.EventID,
		Name:              event.Name,
		Date:              event.Date,
		MSPID:             event.MSPID,
		OrganizerUsername: event.OrganizerUsername,
	}
}

func MapSectionToDTO(section *models.Section, eventID string) *SectionDTO {
	return &SectionDTO{
		EventID:      eventID,
		Name:         section.Name,
		SoldTickets:  section.SoldTickets,
		TotalTickets: section.TotalTickets,
	}
}
