package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/ticken-ticket"
)

func main() {
	// add metadata and init transaction context
	tickenTicketContract := new(ticken_ticket.Contract)
	tickenTicketContract.Info.Version = "0.0.1"
	tickenTicketContract.Name = "ticken-ticket-contract"
	tickenTicketContract.TransactionContextHandler = new(ticken_ticket.TransactionContext)

	cc, err := contractapi.NewChaincode(tickenTicketContract)
	if err != nil {
		panic(err.Error())
	}

	err = cc.Start()
	if err != nil {
		panic(err.Error())
	}
}
