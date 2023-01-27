package common

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ITickenTxContext interface {
	contractapi.TransactionContextInterface
	GetInvoker(chaincode string) *Invoker
	GetContextIdentity() (string, string, error)
}

type TickenTxContext struct {
	contractapi.TransactionContext
}

func NewTransactionContext() *TickenTxContext {
	ctx := new(TickenTxContext)
	return ctx
}

func (ctx *TickenTxContext) GetInvoker(chaincode string) *Invoker {
	return NewInvoker(chaincode, ctx.GetStub())
}

func (ctx *TickenTxContext) GetContextIdentity() (string, string, error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", "", err
	}

	x509Cert, err := ctx.GetClientIdentity().GetX509Certificate()
	if err != nil {
		return "", "", err
	}

	username := x509Cert.Subject.OrganizationalUnit[0]
	return mspID, username, nil
}
