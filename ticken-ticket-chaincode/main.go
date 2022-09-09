package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/tickenticket"
)

func main() {
	// add metadata and init transaction context
	tickenTicketContract := new(tickenticket.Contract)
	tickenTicketContract.Info.Version = "0.0.1"
	tickenTicketContract.Name = "tickenticket-contract"
	tickenTicketContract.TransactionContextHandler = tickenticket.NewTransactionContext()

	cc, err := contractapi.NewChaincode(tickenTicketContract)
	if err != nil {
		panic(err.Error())
	}

	err = cc.Start()
	if err != nil {
		panic(err.Error())
	}
}
