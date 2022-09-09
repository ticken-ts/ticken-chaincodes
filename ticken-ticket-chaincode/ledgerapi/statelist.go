package ledgerapi

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type stateList struct {
	Stub        shim.ChaincodeStubInterface
	Name        string
	Deserialize func([]byte, State) error
}

func NewStateList(stub shim.ChaincodeStubInterface, name string, deserializeFunc func([]byte, State) error) *stateList {
	return &stateList{
		Stub:        stub,
		Name:        name,
		Deserialize: deserializeFunc,
	}
}

func (sl *stateList) AddState(state State) error {
	splitKey := SplitKey(state.GetKey())
	key, _ := sl.Stub.CreateCompositeKey(sl.Name, splitKey)
	data, err := state.Serialize()

	if err != nil {
		return err
	}

	return sl.Stub.PutState(key, data)
}

func (sl *stateList) GetState(key string, state State) error {
	ledgerKey, _ := sl.Stub.CreateCompositeKey(sl.Name, SplitKey(key))
	data, err := sl.Stub.GetState(ledgerKey)

	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("no state found for %s", key)
	}

	return sl.Deserialize(data, state)
}

func (sl *stateList) UpdateState(state State) error {
	return sl.AddState(state)
}
