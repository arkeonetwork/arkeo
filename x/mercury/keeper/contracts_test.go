package keeper

import (
	"mercury/common"
	"mercury/x/mercury/types"

	. "gopkg.in/check.v1"
)

type KeeperContractSuite struct{}

var _ = Suite(&KeeperContractSuite{})

func (s *KeeperContractSuite) TestContract(c *C) {
	ctx, k := SetupKeeper(c)

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

func (s *KeeperContractSuite) TestContractExpirationSet(c *C) {
	var err error
	ctx, k := SetupKeeper(c)
	set := types.ContractExpirationSet{}
	c.Check(k.SetContractExpirationSet(ctx, set), NotNil) // empty asset should error

	set.Height = 100
	c.Check(k.SetContractExpirationSet(ctx, set), IsNil) // empty asset NOT should error

	exp := types.NewContractExpiration(types.GetRandomPubKey(), common.BTCChain, types.GetRandomBech32Addr())
	set.Contracts = append(set.Contracts, exp)

	c.Assert(k.SetContractExpirationSet(ctx, set), IsNil)
	set, err = k.GetContractExpirationSet(ctx, set.Height)
	c.Assert(err, IsNil)
	c.Check(set.Height, Equals, int64(100))
	c.Assert(set.Contracts, HasLen, 1)
	c.Check(set.Contracts[0].ProviderPubKey, Equals, exp.ProviderPubKey)
	c.Check(set.Contracts[0].Chain, Equals, exp.Chain)
	c.Check(set.Contracts[0].ClientAddress.String(), Equals, exp.ClientAddress.String())

	k.RemoveContractExpirationSet(ctx, 100)
	set, err = k.GetContractExpirationSet(ctx, set.Height)
	c.Assert(err, IsNil)
	c.Assert(set.Contracts, HasLen, 0)
}
