package types

import (
	fmt "fmt"
	"mercury/common"
	"mercury/common/cosmos"
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

func NewContract(pubkey common.PubKey, chain common.Chain, client cosmos.AccAddress) Contract {
	return Contract{
		ProviderPubKey: pubkey,
		Chain:          chain,
		ClientAddress:  client,
	}
}

func (c Contract) ExpireHeight() uint64 {
	return c.Height + c.Duration
}

func (c Contract) Key() string {
	return fmt.Sprintf("%s/%s/%s", c.ProviderPubKey, c.Chain, c.ClientAddress)
}
