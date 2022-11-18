package types

import (
	"arkeo/common"
	"arkeo/common/cosmos"
	"encoding/json"
	fmt "fmt"
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

func NewContract(pubkey common.PubKey, chain common.Chain, client common.PubKey) Contract {
	return Contract{
		ProviderPubKey: pubkey,
		Chain:          chain,
		Client:         client,
		Delegate:       common.EmptyPubKey,
		Deposit:        cosmos.ZeroInt(),
		Paid:           cosmos.ZeroInt(),
	}
}

func (c Contract) Key() string {
	return fmt.Sprintf("%s/%s/%s", c.ProviderPubKey, c.Chain, c.FetchSpender())
}

func (c Contract) FetchSpender() common.PubKey {
	if !c.Delegate.IsEmpty() {
		return c.Delegate
	}
	return c.Client
}

func (c Contract) Expiration() int64 {
	return c.Height + c.Duration
}

func (c Contract) IsOpen(h int64) bool {
	if c.IsEmpty() {
		return false
	}
	if c.Expiration() < h {
		return false
	}
	if c.ClosedHeight > 0 && c.ClosedHeight < h {
		return false
	}
	return true
}

func (c Contract) IsClose(h int64) bool {
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

func NewContractExpiration(pubkey common.PubKey, chain common.Chain, client common.PubKey) *ContractExpiration {
	return &ContractExpiration{
		ProviderPubKey: pubkey,
		Chain:          chain,
		Client:         client,
	}
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
