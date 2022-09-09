package tickenticket

type ListInterface interface {
	AddTicket(ticket *Ticket) error
	UpdateTicket(ticket *Ticket) error
	CountTicketsInSection(eventID string, section string) (int, error)
	GetTicket(eventID string, ticketID string) (*Ticket, error)
}
