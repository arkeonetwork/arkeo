package keeper

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

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

func (ManagerSuite) TestContractEndBlock(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)
	mgr := NewManager(k, sk)

	// create a provider for 2 chains
	providerPubKey := types.GetRandomPubKey()
	provider := types.NewProvider(providerPubKey, common.BTCChain)
	provider.Bond = cosmos.NewInt(20000000000)
	provider.LastUpdate = ctx.BlockHeight()
	c.Assert(k.SetProvider(ctx, provider), IsNil)
	provider.Chain = common.ETHChain
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Chain:               common.BTCChain.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
	}
	err := s.ModProviderHandle(ctx, &modProviderMsg)
	c.Assert(err, IsNil)
	modProviderMsg.Chain = common.ETHChain.String()
	err = s.ModProviderHandle(ctx, &modProviderMsg)
	c.Assert(err, IsNil)

	// create user1 to open a contract against the provider.
	user1PubKey := types.GetRandomPubKey()
	user1Address, err := user1PubKey.GetMyAddress()
	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, user1Address, getCoin(common.Tokens(10))), IsNil)

	msg := types.MsgOpenContract{
		Provider:     providerPubKey,
		Chain:        common.BTCChain.String(),
		Creator:      user1Address.String(),
		Client:       user1PubKey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	c.Assert(err, IsNil)

	// have user1 open a contract for a delegate.
	delegatePubKey := types.GetRandomPubKey()
	msg.Delegate = delegatePubKey
	_, err = s.OpenContract(ctx, &msg)
	c.Assert(err, IsNil)

	// create user2 to open a contract against the provider.
	user2PubKey := types.GetRandomPubKey()
	user2Address, err := user2PubKey.GetMyAddress()
	c.Assert(err, IsNil)

	c.Assert(k.MintAndSendToAccount(ctx, user2Address, getCoin(common.Tokens(20))), IsNil)
	msg.Delegate = common.EmptyPubKey
	msg.Client = user2PubKey
	msg.Creator = user2Address.String()
	_, err = s.OpenContract(ctx, &msg)
	c.Assert(err, IsNil)

	// confirm user 1 has an active and open contract
	activeContract, err := k.GetActiveContractForUser(ctx, user1PubKey, providerPubKey, common.BTCChain)
	c.Assert(err, IsNil)
	c.Check(activeContract.IsEmpty(), Equals, false)

	// have user2 open another conrtact with a different expiration
	// to ensure we properly handle a user contract set with multiples
	// contracts with different expiries.
	msg.Duration = 200
	msg.Chain = common.ETHChain.String()
	_, err = s.OpenContract(ctx, &msg)
	c.Assert(err, IsNil)

	// confirm user 2 has 2 contracts in their set.
	contractSet, err := k.GetUserContractSet(ctx, user2PubKey)
	c.Assert(err, IsNil)

	contractIdExpiring := contractSet.ContractSet.ContractIds[0]
	c.Check(len(contractSet.ContractSet.ContractIds), Equals, 2)

	// advance 100 blocks and call end block
	ctx = ctx.WithBlockHeight(110)
	err = mgr.ContractEndBlock(ctx)
	c.Assert(err, IsNil)

	// user 2 should only have 1 contract left in their set.
	contractSet, err = k.GetUserContractSet(ctx, user2PubKey)
	c.Assert(err, IsNil)

	c.Check(len(contractSet.ContractSet.ContractIds), Equals, 1)

	// confirm the contract id left is not the one that expired.
	c.Check(contractIdExpiring, Not(Equals), contractSet.ContractSet.ContractIds[0])

	// cofirm user1 has no active contract.
	activeContract, err = k.GetActiveContractForUser(ctx, user1PubKey, providerPubKey, common.BTCChain)
	c.Assert(err, IsNil)
	c.Check(activeContract.IsEmpty(), Equals, true)

	// advance 100 more blocks and call end block to ensure user 2 has no contracts left.
	ctx = ctx.WithBlockHeight(210)
	err = mgr.ContractEndBlock(ctx)
	c.Assert(err, IsNil)
	contractSet, err = k.GetUserContractSet(ctx, user2PubKey)
	c.Assert(err, IsNil)
	c.Assert(contractSet.ContractSet, IsNil)
}
