package keeper

import (
	"github.com/ArkeoNetwork/arkeo-protocol/common"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/configs"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	. "gopkg.in/check.v1"

	"github.com/cosmos/cosmos-sdk/simapp"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type ManagerSuite struct{}

var _ = Suite(&ManagerSuite{})

func (ManagerSuite) TestValidatorPayout(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)

	pks := simapp.CreateTestPubKeys(3)
	pk1, err := common.NewPubKeyFromCrypto(pks[0])
	c.Assert(err, IsNil)
	acc1, err := pk1.GetMyAddress()
	c.Assert(err, IsNil)
	pk2, err := common.NewPubKeyFromCrypto(pks[1])
	c.Assert(err, IsNil)
	acc2, err := pk2.GetMyAddress()
	c.Assert(err, IsNil)
	pk3, err := common.NewPubKeyFromCrypto(pks[2])
	c.Assert(err, IsNil)
	acc3, err := pk3.GetMyAddress()
	c.Assert(err, IsNil)

	valAddrs := simapp.ConvertAddrsToValAddrs([]cosmos.AccAddress{acc1, acc2, acc3})

	val1, err := stakingtypes.NewValidator(valAddrs[0], pks[0], stakingtypes.Description{})
	c.Assert(err, IsNil)
	val1.Tokens = cosmos.NewInt(100)
	val1.DelegatorShares = cosmos.NewDec(100 + 10 + 20)
	val1.Status = stakingtypes.Bonded
	val1.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(1, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	val2, err := stakingtypes.NewValidator(valAddrs[1], pks[1], stakingtypes.Description{})
	c.Assert(err, IsNil)
	val2.Tokens = cosmos.NewInt(200)
	val2.DelegatorShares = cosmos.NewDec(200 + 20)
	val2.Status = stakingtypes.Bonded
	val2.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(2, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	val3, err := stakingtypes.NewValidator(valAddrs[2], pks[2], stakingtypes.Description{})
	c.Assert(err, IsNil)
	val3.Tokens = cosmos.NewInt(500)
	val3.DelegatorShares = cosmos.NewDec(500)
	val3.Status = stakingtypes.Bonded
	val3.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(5, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	vals := []stakingtypes.Validator{val1, val2, val3}
	for _, val := range vals {
		sk.SetValidator(ctx, val)
		c.Assert(sk.SetValidatorByConsAddr(ctx, val), IsNil)
		sk.SetNewValidatorByPowerIndex(ctx, val)
	}

	delAcc1 := types.GetRandomBech32Addr()
	delAcc2 := types.GetRandomBech32Addr()
	delAcc3 := types.GetRandomBech32Addr()

	sk.SetDelegation(ctx, stakingtypes.NewDelegation(acc1, valAddrs[0], cosmos.NewDec(100)))
	sk.SetDelegation(ctx, stakingtypes.NewDelegation(acc2, valAddrs[1], cosmos.NewDec(200)))
	sk.SetDelegation(ctx, stakingtypes.NewDelegation(acc3, valAddrs[2], cosmos.NewDec(500)))

	del1 := stakingtypes.NewDelegation(delAcc1, valAddrs[0], cosmos.NewDec(10))
	del2 := stakingtypes.NewDelegation(delAcc2, valAddrs[0], cosmos.NewDec(20))
	del3 := stakingtypes.NewDelegation(delAcc3, valAddrs[1], cosmos.NewDec(20))
	sk.SetDelegation(ctx, del1)
	sk.SetDelegation(ctx, del2)
	sk.SetDelegation(ctx, del3)

	c.Assert(k.MintToModule(ctx, types.ModuleName, getCoin(common.Tokens(50000))), IsNil)
	c.Assert(k.SendFromModuleToModule(ctx, types.ModuleName, types.ReserveName, getCoins(common.Tokens(50000))), IsNil)

	mgr := NewManager(k, sk)
	ctx = ctx.WithBlockHeight(mgr.FetchConfig(ctx, configs.ValidatorPayoutCycle))

	blockReward := int64(237792)
	c.Assert(mgr.ValidatorEndBlock(ctx), IsNil)
	c.Check(k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), Equals, 5000000000000-blockReward)

	// check validator balances
	totalBal := cosmos.ZeroInt()
	bal := k.GetBalance(ctx, acc1)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(27984))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, acc2)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(55962))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, acc3)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(139878))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	// check delegate balances
	bal = k.GetBalance(ctx, delAcc1)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(2795))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, delAcc2)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(5589))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, delAcc3)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(5584))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	// ensure block reward is equal to total rewarded to validators and delegates
	c.Check(blockReward, Equals, totalBal.Int64())
}
