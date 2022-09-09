package ccinvokers

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"strings"
	"time"
)

const dateLayout = "1998-09-11T14:36"

type Event struct {
	EventID  string         `json:"id"`
	Name     string         `json:"name"`
	Sections map[string]int `json:"sections"`
	Datetime string         `json:"datetime"`
	Active   bool           `json:"active"`
}

func (e *Event) getSectionCapacity(sectionName string) (int, error) {
	normalizedSectionName := strings.ToUpper(sectionName)

	if capacity, ok := e.Sections[normalizedSectionName]; ok {
		return capacity, nil
	}

	return -1, fmt.Errorf("section %s not found", normalizedSectionName)
}

func (e *Event) ticketSellIsOpen(solicitationTimestamp *timestamp.Timestamp) bool {
	if !e.Active {
		return false
	}

	// we are sure that there is no error in datetime
	// format. it is going to be validated on the other
	// chaincode before inserting
	eventStartTime, _ := time.Parse(e.Datetime, dateLayout)

	solicitationTime := time.Unix(
		solicitationTimestamp.GetSeconds(),
		int64(solicitationTimestamp.GetNanos()),
	)

	// here we are making the supposition that the
	// tickets are available until the start time
	return eventStartTime.Before(solicitationTime)
}
