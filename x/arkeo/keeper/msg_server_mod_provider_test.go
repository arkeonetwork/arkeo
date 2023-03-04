package keeper

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type ModProviderSuite struct{}

var _ = Suite(&ModProviderSuite{})

func (ModProviderSuite) TestValidate(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)

	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()

	provider := types.NewProvider(pubkey, common.BTCChain)
	provider.Bond = cosmos.NewInt(500)
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	// happy path
	msg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Chain:               provider.Chain.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
	}
	c.Assert(s.ModProviderValidate(ctx, &msg), IsNil)

	// bad min duration
	msg.MinContractDuration = 5256000 * 2
	err := s.ModProviderValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInvalidModProviderMinContractDuration)

	// bad max duration
	msg.MinContractDuration = 10
	msg.MaxContractDuration = 5256000 * 2
	err = s.ModProviderValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInvalidModProviderMaxContractDuration)
}

func (ModProviderSuite) TestHandle(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)

	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)

	// happy path
	msg := types.MsgModProvider{
		Creator:             acct.String(),
		Provider:            pubkey,
		Chain:               common.BTCChain.String(),
		MetadataUri:         "foobar",
		MetadataNonce:       3,
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		SubscriptionRate:    11,
		PayAsYouGoRate:      12,
	}
	c.Assert(s.ModProviderHandle(ctx, &msg), IsNil)

	provider, err := k.GetProvider(ctx, msg.Provider, common.BTCChain)
	c.Assert(err, IsNil)
	c.Check(provider.MetadataUri, Equals, "foobar")
	c.Check(provider.MetadataNonce, Equals, uint64(3))
	c.Check(provider.MinContractDuration, Equals, int64(10))
	c.Check(provider.MaxContractDuration, Equals, int64(500))
	c.Check(provider.Status, Equals, types.ProviderStatus_ONLINE)
	c.Check(provider.SubscriptionRate, Equals, int64(11))
	c.Check(provider.PayAsYouGoRate, Equals, int64(12))
}
