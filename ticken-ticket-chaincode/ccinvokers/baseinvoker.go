package ccinvokers

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type baseInvoker struct {
	chaincodeName string
	stub          shim.ChaincodeStubInterface
}

func NewBaseInvoker(chaincodeName string, stub shim.ChaincodeStubInterface) *baseInvoker {
	return &baseInvoker{
		chaincodeName: chaincodeName,
		stub:          stub,
	}
}

func (baseInvoker *baseInvoker) Invoke(opName string, args ...string) ([]byte, error) {
	invokeResponse := baseInvoker.stub.InvokeChaincode(
		baseInvoker.chaincodeName,
		getQueryArgs(opName, args...),
		baseInvoker.stub.GetChannelID(),
	)

	if invokeResponse.Status != shim.OK {
		return nil, fmt.Errorf(invokeResponse.Message)
	}

	return invokeResponse.Payload, nil
}

func getQueryArgs(opName string, args ...string) [][]byte {
	queryArgs := make([][]byte, len(args)+1)

	queryArgs[0] = []byte(opName)
	for i, arg := range args {
		queryArgs[i+1] = []byte(arg)
	}

	return queryArgs
}
