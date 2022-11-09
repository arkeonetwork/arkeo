package keeper

import (
	"arkeo/common"
	"arkeo/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type KeeperProviderSuite struct{}

var _ = Suite(&KeeperProviderSuite{})

func (s *KeeperProviderSuite) TestProvider(c *C) {
	ctx, k := SetupKeeper(c)

	c.Check(k.SetProvider(ctx, types.Provider{}), NotNil) // empty asset should error

	provider := types.NewProvider(types.GetRandomPubKey(), common.BTCChain)

	err := k.SetProvider(ctx, provider)
	c.Assert(err, IsNil)
	provider, err = k.GetProvider(ctx, provider.PubKey, provider.Chain)
	c.Assert(err, IsNil)
	c.Check(provider.Chain.Equals(common.BTCChain), Equals, true)
	c.Check(k.ProviderExists(ctx, provider.PubKey, provider.Chain), Equals, true)
	c.Check(k.ProviderExists(ctx, provider.PubKey, common.ETHChain), Equals, false)

	k.RemoveProvider(ctx, provider.PubKey, provider.Chain)
	c.Check(k.ProviderExists(ctx, provider.PubKey, provider.Chain), Equals, false)
}
