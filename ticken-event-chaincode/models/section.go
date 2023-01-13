package models

import "fmt"

type Section struct {
	Name         string
	TicketPrice  float64
	TotalTickets int
	SoldTickets  int
}

func (section *Section) RemainingTickets() int {
	return section.TotalTickets - section.SoldTickets
}

func (section *Section) IsAvailable() bool {
	return section.RemainingTickets() > 0
}

func (section *Section) SellTicket() error {
	if section.RemainingTickets() == 0 {
		return fmt.Errorf("section %s has no more ticket available", section.Name)
	}

	section.SoldTickets += 1
	return nil
}
