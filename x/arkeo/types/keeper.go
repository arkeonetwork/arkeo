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

func (p Provider) Key() string {
	return fmt.Sprintf("%s/%s", p.PubKey, p.Chain)
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

func (c Contract) Key() string {
	return strconv.FormatUint(c.Id, 10)
}

func (c Contract) GetSpender() common.PubKey {
	if !c.Delegate.IsEmpty() {
		return c.Delegate
	}
	return c.Client
}

func (c Contract) Expiration() int64 {
	return c.Height + c.Duration
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

func (c Contract) IsClosed(h int64) bool {
	return !c.IsOpen(h)
}

func (c Contract) IsEmpty() bool {
	return c.Height == 0
}

func (c Contract) ClientAddress() cosmos.AccAddress {
	addr, err := c.Client.GetMyAddress()
	if err != nil {
		panic(err)
	}
	return addr
}

func (ct *ContractType) UnmarshalJSON(b []byte) error {
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*ct = ContractType(v)
	case string:
		*ct = ContractType(ContractType_value[v])
	}
	return nil
}

func (userContractSet *UserContractSet) RemoveContractFromSet(contractIdToRemove uint64) (*UserContractSet, error) {
	if userContractSet == nil {
		return nil, fmt.Errorf("user contract set is nil")
	}

	if userContractSet.ContractSet == nil {
		return nil, fmt.Errorf("contract set is nil")
	}

	if len(userContractSet.ContractSet.ContractIds) == 0 {
		return nil, fmt.Errorf("contract set is empty")
	}

	isFound := false
	for i, contractId := range userContractSet.ContractSet.ContractIds {
		if contractId == contractIdToRemove {
			userContractSet.ContractSet.ContractIds = append(userContractSet.ContractSet.ContractIds[:i], userContractSet.ContractSet.ContractIds[i+1:]...)
			isFound = true
			break
		}
	}

	if !isFound {
		return userContractSet, fmt.Errorf("contract %d not found in user contract set", contractIdToRemove)
	}
	return userContractSet, nil
}
