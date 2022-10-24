package keeper

import (
	. "gopkg.in/check.v1"

	"mercury/common"
	"mercury/x/mercury/types"
)

type KeeperContractSuite struct{}

var _ = Suite(&KeeperContractSuite{})

func (s *KeeperContractSuite) TestContract(c *C) {
	ctx, k := SetupKeeper()

	c.Check(k.SetContract(ctx, types.Contract{}), NotNil) // empty asset should error

	contract := types.NewContract(types.GetRandomPubKey(), common.BTCChain, types.GetRandomBech32Addr())

	err := k.SetContract(ctx, contract)
	c.Assert(err, IsNil)
	contract, err = k.GetContract(ctx, contract.ProviderPubKey, contract.Chain, contract.ClientAddress)
	c.Assert(err, IsNil)
	c.Check(contract.Chain.Equals(common.BTCChain), Equals, true)
	c.Check(k.ContractExists(ctx, contract.ProviderPubKey, contract.Chain, contract.ClientAddress), Equals, true)
	c.Check(k.ContractExists(ctx, contract.ProviderPubKey, common.ETHChain, contract.ClientAddress), Equals, false)

	k.RemoveContract(ctx, contract.ProviderPubKey, contract.Chain, contract.ClientAddress)
	c.Check(k.ContractExists(ctx, contract.ProviderPubKey, contract.Chain, contract.ClientAddress), Equals, false)
}
