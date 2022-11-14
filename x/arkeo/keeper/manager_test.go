package keeper

import (
	"arkeo/common"
	"arkeo/common/cosmos"
	"arkeo/x/arkeo/configs"
	"arkeo/x/arkeo/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "gopkg.in/check.v1"
)

type ManagerSuite struct{}

var _ = Suite(&ManagerSuite{})

type TestManagerKeeperValidatorPayout struct {
	KVStoreDummy
	keeper Keeper
}

func (k *TestManagerKeeperValidatorPayout) MintToModule(ctx cosmos.Context, module string, coin cosmos.Coin) error {
	return k.keeper.MintToModule(ctx, module, coin)
}

func (k *TestManagerKeeperValidatorPayout) SendFromModuleToModule(ctx cosmos.Context, from, to string, coins cosmos.Coins) error {
	return k.keeper.SendFromModuleToModule(ctx, from, to, coins)
}

func (k *TestManagerKeeperValidatorPayout) SendFromModuleToAccount(ctx cosmos.Context, from string, to cosmos.AccAddress, coins cosmos.Coins) error {
	return k.keeper.SendFromModuleToAccount(ctx, from, to, coins)
}

func (k *TestManagerKeeperValidatorPayout) GetBalanceOfModule(ctx cosmos.Context, moduleName, denom string) cosmos.Int {
	return k.keeper.GetBalanceOfModule(ctx, moduleName, denom)
}

func (k *TestManagerKeeperValidatorPayout) GetBalance(ctx cosmos.Context, acc cosmos.AccAddress) cosmos.Coins {
	return k.keeper.GetBalance(ctx, acc)
}

func (k *TestManagerKeeperValidatorPayout) GetActiveValidators(ctx cosmos.Context) []stakingtypes.Validator {
	return []stakingtypes.Validator{
		{
			OperatorAddress: "cosmos1n07t9h5h5f75aqmfpdv6tlp490dktx2mynkmxw",
			Tokens:          cosmos.NewInt(common.Tokens(100)),
			Status:          stakingtypes.Bonded,
		},
		{
			OperatorAddress: "cosmos1tvvzxkzgwavypl8agjtmlk5tzdtrh7x5scfvmn",
			Tokens:          cosmos.NewInt(common.Tokens(100)),
			Status:          stakingtypes.Bonded,
		},
		{
			OperatorAddress: "cosmos1qwf8gk53r3x93v975ymfsgfgp4m846hj025qug",
			Tokens:          cosmos.NewInt(common.Tokens(200)),
			Status:          stakingtypes.Bonded,
		},
	}
}

func (ManagerSuite) TestValidatorPayout(c *C) {
	ctx, keepr, sk := SetupKeeperWithStaking(c)
	k := &TestManagerKeeperValidatorPayout{
		keeper: keepr,
	}

	c.Assert(k.MintToModule(ctx, types.ModuleName, getCoin(common.Tokens(500))), IsNil)
	c.Assert(k.SendFromModuleToModule(ctx, types.ModuleName, types.ReserveName, getCoins(common.Tokens(500))), IsNil)

	mgr := NewManager(k, sk)
	ctx = ctx.WithBlockHeight(mgr.FetchConfig(ctx, configs.ValidatorPayoutCycle))

	c.Assert(mgr.ValidatorEndBlock(ctx), IsNil)
	c.Assert(k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), Equals, int64(49998573223))

	validators := k.GetActiveValidators(ctx)
	c.Assert(validators, HasLen, 3)

	valAddr, err := cosmos.AccAddressFromBech32("cosmos1n07t9h5h5f75aqmfpdv6tlp490dktx2mynkmxw")
	c.Assert(err, IsNil)
	bal := k.GetBalance(ctx, valAddr)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(356694))

	valAddr, err = cosmos.AccAddressFromBech32("cosmos1tvvzxkzgwavypl8agjtmlk5tzdtrh7x5scfvmn")
	c.Assert(err, IsNil)
	bal = k.GetBalance(ctx, valAddr)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(356694))

	valAddr, err = cosmos.AccAddressFromBech32("cosmos1qwf8gk53r3x93v975ymfsgfgp4m846hj025qug")
	c.Assert(err, IsNil)
	bal = k.GetBalance(ctx, valAddr)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(713389))
}
