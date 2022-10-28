package ccnotifier

import (
	"ticken-event-contract/models"
	"time"
)

type SectionDTO struct {
	Name             string `json:"name"`
	TotalTickets     int    `json:"total_tickets"`
	RemainingTickets int    `json:"remaining_tickets"`
}

type EventDTO struct {
	EventID        string       `json:"event_id"`
	Name           string       `json:"name"`
	Date           time.Time    `json:"date"`
	Sections       []SectionDTO `json:"sections"`
	OrganizationID string       `json:"organization_id"`
}

func MapEventToNotificationDTO(event *models.Event) *EventDTO {
	sectionsDTO := make([]SectionDTO, len(event.Sections))
	for i, section := range event.Sections {
		sectionsDTO[i] = SectionDTO{
			Name:             section.Name,
			RemainingTickets: section.RemainingTickets,
			TotalTickets:     section.TotalTickets,
		}
	}

	return &EventDTO{
		EventID:        event.EventID,
		Name:           event.Name,
		Date:           event.Date,
		Sections:       sectionsDTO,
		OrganizationID: event.OrganizationID,
	}
}
