package ccnotifier

import "ticken-event-contract/models"

type Notifier interface {
	NotifyEventCreation(event *models.Event) error
	NotifySectionAddition(event *models.Event) error
}
