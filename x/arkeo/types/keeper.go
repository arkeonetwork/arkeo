package types

import (
	"encoding/json"
	fmt "fmt"
	"strconv"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
)

func NewProvider(pubkey common.PubKey, chain common.Chain) Provider {
	return Provider{
		PubKey: pubkey,
		Chain:  chain,
		Bond:   cosmos.ZeroInt(),
	}
}

func (provider Provider) Key() string {
	return fmt.Sprintf("%s/%s", provider.PubKey, provider.Chain)
}

func NewContract(provider common.PubKey, chain common.Chain, client common.PubKey) Contract {
	return Contract{
		Provider: provider,
		Chain:    chain,
		Client:   client,
		Delegate: common.EmptyPubKey,
		Deposit:  cosmos.ZeroInt(),
		Paid:     cosmos.ZeroInt(),
	}
}

func (contract Contract) Key() string {
	return strconv.FormatUint(contract.Id, 10)
}

func (contract Contract) GetSpender() common.PubKey {
	if !contract.Delegate.IsEmpty() {
		return contract.Delegate
	}
	return contract.Client
}

func (contract Contract) Expiration() int64 {
	return contract.Height + contract.Duration
}

func (c Contract) IsOpen(height int64) bool {
	if c.IsEmpty() {
		return false
	}
	if c.Expiration() < height {
		return false
	}
	if c.ClosedHeight > 0 && c.ClosedHeight < height {
		return false
	}
	return true
}

func (contract Contract) IsClosed(h int64) bool {
	return !contract.IsOpen(h)
}

func (contract Contract) IsEmpty() bool {
	return contract.Height == 0
}

func (contract Contract) ClientAddress() cosmos.AccAddress {
	addr, err := contract.Client.GetMyAddress()
	if err != nil {
		panic(err)
	}
	return addr
}

func (contractType *ContractType) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*contractType = ContractType(v)
	case string:
		*contractType = ContractType(ContractType_value[v])
	}
	return nil
}

func (userContractSet *UserContractSet) RemoveContractFromSet(contractIdToRemove uint64) error {
	if userContractSet == nil {
		return fmt.Errorf("user contract set is nil")
	}

	if userContractSet.ContractSet == nil {
		return fmt.Errorf("contract set is nil")
	}

	if len(userContractSet.ContractSet.ContractIds) == 0 {
		return fmt.Errorf("contract set is empty")
	}

	for i, contractId := range userContractSet.ContractSet.ContractIds {
		if contractId == contractIdToRemove {
			userContractSet.ContractSet.ContractIds = append(userContractSet.ContractSet.ContractIds[:i], userContractSet.ContractSet.ContractIds[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("contract %d not found in user contract set", contractIdToRemove)
}
