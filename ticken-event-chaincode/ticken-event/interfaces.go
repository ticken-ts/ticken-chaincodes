package ticken_event

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-event-contract/ccnotifier"
	"ticken-event-contract/models"
)

type EventListInterface interface {
	AddEvent(event *models.Event) error
	UpdateEvent(event *models.Event) error
	EventExists(event *models.Event) (bool, error)
	GetEvent(evenID string) (*models.Event, error)
}

type ITickenTxContext interface {
	contractapi.TransactionContextInterface
	GetEventList() EventListInterface
	GetNotifier() ccnotifier.Notifier
}
