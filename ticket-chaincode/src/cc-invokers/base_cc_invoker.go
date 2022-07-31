package cc_invokers

import (
	"fmt"
	ticken_ticket "ticken-ticket-contract/ticken-ticket"
)

const (
	SuccessChaincodeInvokeResult = 200
)

type CCBaseInvoker interface {
	Invoke(opName string, args ...string) ([]byte, error)
}

type ccBaseInvoker struct {
	ccName string
	ctx    ticken_ticket.TransactionContextInterface
}

func NewCCBaseInvoker(ccName string, ctx ticken_ticket.TransactionContextInterface) *ccBaseInvoker {
	ccBaseInvoker := new(ccBaseInvoker)
	ccBaseInvoker.ccName = ccName
	ccBaseInvoker.ctx = ctx
	return ccBaseInvoker
}

func (i *ccBaseInvoker) Invoke(opName string, args ...string) ([]byte, error) {
	invokeReponse := i.ctx.GetStub().InvokeChaincode(
		i.ccName,
		getQueryArgs(opName, args...),
		i.ctx.GetStub().GetChannelID(),
	)

	if invokeReponse.Status != SuccessChaincodeInvokeResult {
		return nil, fmt.Errorf(invokeReponse.Message)
	}

	return invokeReponse.Payload, nil
}

func getQueryArgs(opName string, args ...string) [][]byte {
	queryArgs := make([][]byte, len(args)+1)

	queryArgs[0] = []byte(opName)
	for i, arg := range args {
		queryArgs[i+1] = []byte(arg)
	}

	return queryArgs
}
