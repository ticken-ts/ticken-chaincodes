package ccinvokers

type BaseInvoker interface {
	Invoke(opName string, args ...string) ([]byte, error)
}

type TickenEventInvoker interface {
	EventExists(eventID string) (bool, error)
	IsAvailable(eventID string, section string) (bool, error)
	AddTicket(eventID string, section string) error
}
