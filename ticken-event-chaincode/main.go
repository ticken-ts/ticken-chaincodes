package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"ticken-ticket-contract/ticken-event"
)

func main() {
	// add metadata and init transaction context
	tickenEventContract := new(ticken_event.Contract)
	tickenEventContract.Info.Version = "0.0.1"
	tickenEventContract.Name = "ticken-event-contract"
	tickenEventContract.TransactionContextHandler = ticken_event.NewTransactionContext()

	cc, err := contractapi.NewChaincode(tickenEventContract)
	if err != nil {
		panic(err.Error())
	}

	err = cc.Start()
	if err != nil {
		panic(err.Error())
	}
}
