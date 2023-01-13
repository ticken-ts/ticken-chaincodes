package ccnotifier

import (
	"ticken-event-contract/models"
	"time"
)

type SectionDTO struct {
	Name             string `json:"name"`
	TotalTickets     int    `json:"total_tickets"`
	RemainingTickets int    `json:"remaining_tickets"`
	EventID          string `json:"event_id"`
}

type EventDTO struct {
	EventID        string    `json:"event_id"`
	Name           string    `json:"name"`
	Date           time.Time `json:"date"`
	OrganizationID string    `json:"organization_id"`
}

func MapEventToDTO(event *models.Event) *EventDTO {
	return &EventDTO{
		EventID:        event.EventID,
		Name:           event.Name,
		Date:           event.Date,
		OrganizationID: event.OrganizationID,
	}
}

func MapSectionToDTO(section *models.Section, eventID string) *SectionDTO {
	return &SectionDTO{
		Name:             section.Name,
		RemainingTickets: section.RemainingTickets,
		TotalTickets:     section.TotalTickets,
		EventID:          eventID,
	}
}
