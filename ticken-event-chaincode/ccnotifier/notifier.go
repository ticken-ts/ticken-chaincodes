package ccnotifier

import (
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"ticken-event-contract/models"
)

type NotificationType string

const (
	EventCreatedNotification NotificationType = "event-created"
	SectionAddedNotification NotificationType = "section-added"
)

type CCNotifier struct {
	stub shim.ChaincodeStubInterface
}

func NewNotifier(stub shim.ChaincodeStubInterface) *CCNotifier {
	return &CCNotifier{stub: stub}
}

func (notifier *CCNotifier) NotifyEventCreation(event *models.Event) error {
	eventDTO := MapEventToDTO(event)
	bytes, err := json.Marshal(eventDTO)

	if err != nil {
		return err
	}

	return notifier.Notify(bytes, EventCreatedNotification)
}

func (notifier *CCNotifier) NotifySectionAddition(section *models.Section, eventID string) error {
	sectionDTO := MapSectionToDTO(section, eventID)
	bytes, err := json.Marshal(sectionDTO)
	if err != nil {
		return err
	}

	return notifier.Notify(bytes, SectionAddedNotification)
}

func (notifier *CCNotifier) Notify(data []byte, notificationType NotificationType) error {
	if err := notifier.stub.SetEvent(string(notificationType), data); err != nil {
		return err
	}
	return nil
}
