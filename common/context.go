package common

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ITickenTxContext interface {
	contractapi.TransactionContextInterface
	GetInvoker(chaincode string) *Invoker
}

type TickenTxContext struct {
	contractapi.TransactionContext
	invokers map[string]*Invoker
}

func NewTransactionContext() *TickenTxContext {
	ctx := new(TickenTxContext)
	ctx.invokers = make(map[string]*Invoker)
	return ctx
}

func (ctx *TickenTxContext) GetInvoker(chaincode string) *Invoker {
	invoker, ok := ctx.invokers[chaincode]
	if !ok {
		invoker = NewInvoker(chaincode, ctx.GetStub())
		ctx.invokers[chaincode] = invoker
	}
	return invoker
}
