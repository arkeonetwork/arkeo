package keeper

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type KeeperContractSuite struct{}

var _ = Suite(&KeeperContractSuite{})

func (s *KeeperContractSuite) TestContract(c *C) {
	ctx, k := SetupKeeper(c)

	c.Check(k.SetContract(ctx, types.Contract{}), NotNil) // empty asset should error

	contract := types.NewContract(types.GetRandomPubKey(), common.BTCChain, types.GetRandomPubKey())
	contract.Id = 1
	err := k.SetContract(ctx, contract)
	c.Assert(err, IsNil)

	contract, err = k.GetContract(ctx, contract.Id)
	c.Assert(err, IsNil)
	c.Check(contract.Chain.Equals(common.BTCChain), Equals, true)
	c.Check(k.ContractExists(ctx, contract.Id), Equals, true)
	c.Check(k.ContractExists(ctx, contract.Id+1), Equals, false)

	k.RemoveContract(ctx, contract.Id)
	c.Check(k.ContractExists(ctx, contract.Id), Equals, false)
}

func (s *KeeperContractSuite) TestContractExpirationSet(c *C) {
	var err error
	ctx, k := SetupKeeper(c)
	set := types.ContractExpirationSet{}
	c.Check(k.SetContractExpirationSet(ctx, set), NotNil) // empty asset should error

	set.Height = 100
	set.ContractSet = &types.ContractSet{}
	c.Check(k.SetContractExpirationSet(ctx, set), IsNil) // empty asset NOT should error

	contractId := uint64(100)
	set.ContractSet.ContractIds = append(set.ContractSet.ContractIds, contractId)

	c.Assert(k.SetContractExpirationSet(ctx, set), IsNil)
	set, err = k.GetContractExpirationSet(ctx, set.Height)
	c.Assert(err, IsNil)
	c.Check(set.Height, Equals, int64(100))
	c.Assert(set.ContractSet.ContractIds, HasLen, 1)

	k.RemoveContractExpirationSet(ctx, 100)
	set, err = k.GetContractExpirationSet(ctx, set.Height)
	c.Assert(err, IsNil)
	c.Assert(set.ContractSet, IsNil)
}
