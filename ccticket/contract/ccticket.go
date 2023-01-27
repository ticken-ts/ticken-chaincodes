package contract

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/ticken-ts/ticken-chaincodes/common"
)

type Contract struct {
	contractapi.Contract
}

const Name = "cc-ticket"

// ********** cc-event integration ********** //

const ccEventName = "cc-event"
const ccEventSellTicketFunc = "SellTicket"

// *****+************************************ //

type Ticket struct {
	TicketID string `json:"ticket_id"`
	Status   Status `json:"status"`

	EventID string `json:"event_id"`
	Section string `json:"section"`

	// represents the owner id in the
	// web service database
	OwnerID string `json:"owner"`
}

type Status string

const (
	ISSUED Status = "issued"
)

func (c *Contract) Issue(ctx common.ITickenTxContext, ticketID, eventID, section, ownerID string) (*Ticket, error) {
	existentTicket, err := c.GetTicket(ctx, ticketID)
	if existentTicket != nil {
		return nil, ccErr("ticket with ID %s already exists", ticketID)
	}

	ownerIDParsed, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, ccErr("error parsing owner id: %v", err)
	}
	eventIDParsed, err := uuid.Parse(eventID)
	if err != nil {
		return nil, ccErr("error parsing event id: %v", err)
	}
	ticketIDParsed, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, ccErr("error parsing ticket id: %v", err)
	}

	ticket := Ticket{
		TicketID: ticketIDParsed.String(),

		EventID: eventIDParsed.String(),
		Section: section,
		Status:  ISSUED,

		OwnerID: ownerIDParsed.String(),
	}

	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		return nil, ccErr("failed to serialize ticket: %v", err)
	}

	// add ticket into the chaincode cc-event
	// note: this operation is atomically handled
	// by the orderers. So, the ticket and the ticket
	// count are updated simultaneously in the same tx
	ccEventSellTicketResponse := ctx.GetStub().InvokeChaincode(
		ccEventName,
		getCCCallArgs(ccEventSellTicketFunc, eventID, section),
		ctx.GetStub().GetChannelID(),
	)

	if ccEventSellTicketResponse.Status != shim.OK {
		return nil, ccErr(ccEventSellTicketResponse.Message)
	}

	if err := ctx.GetStub().PutState(ticketID, ticketJSON); err != nil {
		return nil, ccErr("failed to updated the state: %v", err)
	}

	return &ticket, nil
}

func (c *Contract) GetTicket(ctx contractapi.TransactionContextInterface, ticketID string) (*Ticket, error) {
	ticketJSON, err := ctx.GetStub().GetState(ticketID)
	if err != nil {
		return nil, ccErr("failed to read ticket: %v", err)
	}
	if ticketJSON == nil {
		return nil, ccErr("ticket %s does not exist", ticketID)
	}

	var ticket Ticket
	if err := json.Unmarshal(ticketJSON, &ticket); err != nil {
		return nil, ccErr("failed to deserialize ticket: %v", err)
	}

	return &ticket, err
}

func ccErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args)
	return fmt.Errorf("[%s] | %v", Name, msg)
}

func getCCCallArgs(opName string, args ...string) [][]byte {
	queryArgs := make([][]byte, len(args)+1)

	queryArgs[0] = []byte(opName)
	for i, arg := range args {
		queryArgs[i+1] = []byte(arg)
	}

	return queryArgs
}
