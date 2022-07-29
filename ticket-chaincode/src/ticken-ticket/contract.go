package ticken_ticket

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Event struct {
	EventID string `json:"event_id"`
	Name    string `json:"name"`
}

const (
	SUCCESS_CHAINCODE_INVOKE_RESULT = 200
	ERROR_CHAINCODE_INVOKE_RESULT   = 500
)

type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

func (c *Contract) Issue(ctx TransactionContextInterface, ticketID string, eventID string, owner string) (*Ticket, error) {
	ticketList := ctx.GetTicketList()

	ticket := new(Ticket)
	ticket.Init(ticketID, eventID, owner)

	params := []string{"get", eventID}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	res := ctx.GetStub().InvokeChaincode("ticken-event", queryArgs, ctx.GetStub().GetChannelID())
	if res.Status != SUCCESS_CHAINCODE_INVOKE_RESULT {
		return nil, fmt.Errorf(res.Message)
	}

	var event Event
	err := json.Unmarshal(res.Payload, &event)
	if err != nil {
		return nil, err
	}

	println("Event found - ID: %s", event.EventID)

	err = ticketList.AddTicket(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (c *Contract) Sign(ctx TransactionContextInterface, eventName string, ticketID string) (*Ticket, error) {
	ticket, err := ctx.GetTicketList().GetTicket(eventName, ticketID)
	if err != nil {
		return nil, err
	}

	ticket.Sign()
	err = ctx.GetTicketList().UpdateTicket(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (c *Contract) Scan(ctx TransactionContextInterface, eventName string, ticketID string) (*Ticket, error) {
	ticket, err := ctx.GetTicketList().GetTicket(eventName, ticketID)
	if err != nil {
		return nil, err
	}

	ticket.Scan()
	err = ctx.GetTicketList().UpdateTicket(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
