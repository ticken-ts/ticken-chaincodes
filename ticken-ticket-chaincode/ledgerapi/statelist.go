package ledgerapi

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// StateList useful for managing putting data in and out
// of the ledger. Implementation of StateListInterface
type StateList struct {
	Ctx         contractapi.TransactionContextInterface
	Name        string
	Deserialize func([]byte, State) error
}

func (sl *StateList) AddState(state State) error {
	splitKey := SplitKey(state.GetKey())
	key, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, splitKey)
	data, err := state.Serialize()

	if err != nil {
		return err
	}

	return sl.Ctx.GetStub().PutState(key, data)
}

func (sl *StateList) GetState(key string, state State) error {
	ledgerKey, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, SplitKey(key))
	data, err := sl.Ctx.GetStub().GetState(ledgerKey)

	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("no state found for %s", key)
	}

	return sl.Deserialize(data, state)
}

func (sl *StateList) UpdateState(state State) error {
	return sl.AddState(state)
}
