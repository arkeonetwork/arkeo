package keeper

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestContractEndBlock(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)
	mgr := NewManager(k, sk)
	creatorAddress := types.GetRandomBech32Addr()

	// create a provider for 2 services
	providerPubKey := types.GetRandomPubKey()
	provider := types.NewProvider(providerPubKey, common.BTCService)
	provider.Bond = cosmos.NewInt(20000000000)
	provider.LastUpdate = ctx.BlockHeight()
	require.NoError(t, k.SetProvider(ctx, provider))
	provider.Service = common.ETHService
	require.NoError(t, k.SetProvider(ctx, provider))

	rates, err := cosmos.ParseCoins("15uarkeo")
	require.NoError(t, err)

	modProviderMsg := types.MsgModProvider{
		Creator:             creatorAddress.String(),
		Provider:            provider.PubKey,
		Service:             common.BTCService.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      rates,
		SubscriptionRate:    rates,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)
	require.NoError(t, err)
	modProviderMsg.Service = common.ETHService.String()
	err = s.ModProviderHandle(ctx, &modProviderMsg)
	require.NoError(t, err)

	// create user1 to open a contract against the provider.
	user1PubKey := types.GetRandomPubKey()
	user1Address, err := user1PubKey.GetMyAddress()
	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, user1Address, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:     providerPubKey.String(),
		Service:      common.BTCService.String(),
		Creator:      user1Address.String(),
		Client:       user1PubKey.String(),
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         rates[0],
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// have user1 open a contract for a delegate.
	delegatePubKey := types.GetRandomPubKey()
	msg.Delegate = delegatePubKey.String()
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// create user2 to open a contract against the provider.
	user2PubKey := types.GetRandomPubKey()
	user2Address, err := user2PubKey.GetMyAddress()
	require.NoError(t, err)

	require.NoError(t, k.MintAndSendToAccount(ctx, user2Address, getCoin(common.Tokens(20))))
	msg.Delegate = common.EmptyPubKey.String()
	msg.Client = user2PubKey.String()
	msg.Creator = user2Address.String()
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// confirm user 1 has an active and open contract
	activeContract, err := k.GetActiveContractForUser(ctx, user1PubKey, providerPubKey, common.BTCService)
	require.NoError(t, err)
	require.False(t, activeContract.IsEmpty())

	// have user2 open another conrtact with a different expiration
	// to ensure we properly handle a user contract set with multiples
	// contracts with different expiries.
	msg.Duration = 200
	msg.Service = common.ETHService.String()
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// confirm user 2 has 2 contracts in their set.
	contractSet, err := k.GetUserContractSet(ctx, user2PubKey)
	require.NoError(t, err)

	contractIdExpiring := contractSet.ContractSet.ContractIds[0]
	require.Len(t, contractSet.ContractSet.ContractIds, 2)

	// advance 100 blocks and call end block
	ctx = ctx.WithBlockHeight(110)
	err = mgr.ContractEndBlock(ctx)
	require.NoError(t, err)

	// user 2 should only have 1 contract left in their set.
	contractSet, err = k.GetUserContractSet(ctx, user2PubKey)
	require.NoError(t, err)

	require.Len(t, contractSet.ContractSet.ContractIds, 1)

	// confirm the contract id left is not the one that expired.
	require.NotEqual(t, contractIdExpiring, contractSet.ContractSet.ContractIds[0])

	// cofirm user1 has no active contract.
	activeContract, err = k.GetActiveContractForUser(ctx, user1PubKey, providerPubKey, common.BTCService)
	require.NoError(t, err)
	require.True(t, activeContract.IsEmpty())

	// advance 100 more blocks and call end block to ensure user 2 has no contracts left.
	ctx = ctx.WithBlockHeight(210)
	err = mgr.ContractEndBlock(ctx)
	require.NoError(t, err)
	contractSet, err = k.GetUserContractSet(ctx, user2PubKey)
	require.NoError(t, err)
	require.Nil(t, contractSet.ContractSet)
}

func TestContractEndBlockWithSettlementDuration(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)
	mgr := NewManager(k, sk)

	creatorAddress := types.GetRandomBech32Addr()

	// create a provider for 2 services
	providerPubKey := types.GetRandomPubKey()
	provider := types.NewProvider(providerPubKey, common.BTCService)
	provider.Bond = cosmos.NewInt(20000000000)
	provider.LastUpdate = ctx.BlockHeight()
	provider.SettlementDuration = 10
	require.NoError(t, k.SetProvider(ctx, provider))
	provider.Service = common.ETHService
	require.NoError(t, k.SetProvider(ctx, provider))

	rates, err := cosmos.ParseCoins("15uarkeo")
	require.NoError(t, err)

	modProviderMsg := types.MsgModProvider{
		Creator:             creatorAddress.String(),
		Provider:            provider.PubKey,
		Service:             common.BTCService.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      rates,
		SubscriptionRate:    rates,
		SettlementDuration:  10,
	}

	err = s.ModProviderHandle(ctx, &modProviderMsg)
	require.NoError(t, err)
	modProviderMsg.Service = common.ETHService.String()
	err = s.ModProviderHandle(ctx, &modProviderMsg)
	require.NoError(t, err)

	// create user1 to open a contract against the provider.
	user1PubKey := types.GetRandomPubKey()
	user1Address, err := user1PubKey.GetMyAddress()
	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, user1Address, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:           providerPubKey.String(),
		Service:            common.BTCService.String(),
		Creator:            user1Address.String(),
		Client:             user1PubKey.String(),
		ContractType:       types.ContractType_PAY_AS_YOU_GO,
		Duration:           100,
		Rate:               rates[0],
		Deposit:            cosmos.NewInt(1500),
		SettlementDuration: 10,
	}
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// get the active contract for user 1
	activeContract, err := k.GetActiveContractForUser(ctx, user1PubKey, providerPubKey, common.BTCService)
	require.NoError(t, err)
	require.False(t, activeContract.IsEmpty())
	require.True(t, activeContract.IsOpen(ctx.BlockHeight()))

	// advance 100 blocks and call end block
	ctx = ctx.WithBlockHeight(111)
	require.True(t, activeContract.IsExpired(ctx.BlockHeight()))

	// call end block which shouldn't yet do anything as the settlement duration hasn't been reached
	err = mgr.ContractEndBlock(ctx)
	require.NoError(t, err)

	activeContract, err = k.GetContract(ctx, activeContract.Id)
	require.NoError(t, err)
	require.Equal(t, activeContract.SettlementHeight, int64(0))

	// advance 10 more blocks and call end block
	ctx = ctx.WithBlockHeight(activeContract.SettlementPeriodEnd())
	err = mgr.ContractEndBlock(ctx)
	require.NoError(t, err)

	// the end block should have settled the contract and set the settlement height
	activeContract, err = k.GetContract(ctx, activeContract.Id)
	require.NoError(t, err)
	require.Equal(t, activeContract.SettlementHeight, activeContract.SettlementPeriodEnd())
}

func TestInvariantBondModule(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	mgr := NewManager(k, sk)

	require.NoError(t, mgr.invariantBondModule(ctx))

	// setup provider
	pubkey := types.GetRandomPubKey()
	provider := types.NewProvider(pubkey, common.BTCService)
	provider.Bond = cosmos.NewInt(500)
	require.NoError(t, k.SetProvider(ctx, provider))

	// invariant should not fire
	require.ErrorIs(t, mgr.invariantBondModule(ctx), types.ErrInvariantBondModule)

	// mint tokens into the provider module, and check that the invariant no longer fires
	require.NoError(t, k.MintToModule(ctx, types.ModuleName, getCoin(1000)))
	require.NoError(t, k.SendFromModuleToModule(ctx, types.ModuleName, types.ProviderName, getCoins(1000)))
	require.NoError(t, mgr.invariantBondModule(ctx))
}

func TestInvariantContractModule(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	mgr := NewManager(k, sk)

	require.NoError(t, mgr.invariantContractModule(ctx))

	// setup provider
	pubkey := types.GetRandomPubKey()
	contract := types.NewContract(pubkey, common.BTCService, pubkey)
	contract.Rate = cosmos.NewInt64Coin(configs.Denom, 10)
	contract.Deposit = cosmos.NewInt(500)
	contract.Paid = cosmos.NewInt(200)
	contract.Duration = 10
	require.NoError(t, k.SetContract(ctx, contract))

	// invariant should fire
	require.ErrorIs(t, mgr.invariantContractModule(ctx), types.ErrInvariantContractModule)

	// mint tokens into the provider module, and check that the invariant no longer fires
	require.NoError(t, k.MintToModule(ctx, types.ModuleName, getCoin(1000)))
	require.NoError(t, k.SendFromModuleToModule(ctx, types.ModuleName, types.ContractName, getCoins(1000)))
	require.NoError(t, mgr.invariantContractModule(ctx))
}

func TestInvariantMaxSupply(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	mgr := NewManager(k, sk)

	require.NoError(t, mgr.invariantMaxSupply(ctx))

	// mint many coins and trigger the invariant
	require.NoError(t, k.MintToModule(ctx, types.ModuleName, getCoin(200_000_000*1e8)))
	require.ErrorIs(t, mgr.invariantMaxSupply(ctx), types.ErrInvariantMaxSupply)
}

func TestParamsRewardsPercentage(t *testing.T) {
	ctx, k, _ := SetupKeeperWithStaking(t)

	params := k.GetParams(ctx)

	require.Equal(t, params.CommunityPoolPercentage.Int64(), int64(10))
}
func TestCommunityPoolDistributionToFoundationCommunityPool(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	mgr := NewManager(k, sk)

	require.NoError(t, k.MintToModule(ctx, disttypes.ModuleName, getCoin(common.Tokens(200000))))

	err := mgr.keeper.MoveTokensFromDistributionToFoundationPoolAccount(ctx)
	require.NoError(t, err)

	address, err := sdk.AccAddressFromBech32(types.FoundationCommunityAccount)
	require.NoError(t, err)

	bal := mgr.keeper.GetBalance(ctx, address).AmountOf(configs.Denom)

	require.Equal(t, bal, sdkmath.NewInt(10000))
}

func TestBlockRewardCalculation(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	mgr := NewManager(k, sk)

	// BlockPer Year -> 5000
	// Emission Curve -> 10
	// Total Reserve -> 100000000
	// validator cycle -> 100
	// reward = (totalReserve / emissionCurve) / (blocksPerYear / valCycle)) -> 2000
	valCycle := int64(100)
	emissionCurve := int64(10)
	blocksPerYear := int64(5000)
	totalReserve := int64(1000000)

	reward := mgr.calcBlockReward(ctx, totalReserve, emissionCurve, blocksPerYear, valCycle)

	require.Equal(t, reward.Amount.Int64(), int64(2000))

	valCycle = int64(10)
	emissionCurve = int64(5)
	blocksPerYear = int64(200)
	totalReserve = int64(999999)

	reward = mgr.calcBlockReward(ctx, totalReserve, emissionCurve, blocksPerYear, valCycle)

	require.Equal(t, reward.Amount.Int64(), int64(10000)) // its 9999.99 rounded to 10000
}

func TestValidatorPayouts(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	mgr := NewManager(k, sk)

	valCycle := int64(100)
	emissionCurve := int64(10)
	blocksPerYear := int64(5000)
	totalReserve := int64(1000000000)

	blockReward := mgr.calcBlockReward(ctx, totalReserve, emissionCurve, blocksPerYear, valCycle)

	require.Equal(t, blockReward.Amount.Int64(), int64(2000000))

	pks := simtestutil.CreateTestPubKeys(3)
	pk1, err := common.NewPubKeyFromCrypto(pks[0])
	require.NoError(t, err)
	acc1, err := pk1.GetMyAddress()
	require.NoError(t, err)
	pk2, err := common.NewPubKeyFromCrypto(pks[1])
	require.NoError(t, err)
	acc2, err := pk2.GetMyAddress()
	require.NoError(t, err)
	pk3, err := common.NewPubKeyFromCrypto(pks[2])
	require.NoError(t, err)
	acc3, err := pk3.GetMyAddress()
	require.NoError(t, err)

	valAddrs := simtestutil.ConvertAddrsToValAddrs([]cosmos.AccAddress{acc1, acc2, acc3})

	val1, err := stakingtypes.NewValidator(valAddrs[0].String(), pks[0], stakingtypes.Description{})
	require.NoError(t, err)
	val1.Tokens = cosmos.NewInt(100)
	val1.DelegatorShares = cosmos.NewDec(130) // Validator + Delegations
	val1.Status = stakingtypes.Bonded
	val1.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(1, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	val2, err := stakingtypes.NewValidator(valAddrs[1].String(), pks[1], stakingtypes.Description{})
	require.NoError(t, err)
	val2.Tokens = cosmos.NewInt(200)
	val2.DelegatorShares = cosmos.NewDec(220)
	val2.Status = stakingtypes.Bonded
	val2.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(2, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	val3, err := stakingtypes.NewValidator(valAddrs[2].String(), pks[2], stakingtypes.Description{})
	require.NoError(t, err)
	val3.Tokens = cosmos.NewInt(500)
	val3.DelegatorShares = cosmos.NewDec(500)
	val3.Status = stakingtypes.Bonded
	val3.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(5, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	vals := []stakingtypes.Validator{val1, val2, val3}
	for _, val := range vals {
		require.NoError(t, sk.SetValidator(ctx, val))
		require.NoError(t, sk.SetValidatorByConsAddr(ctx, val))
		require.NoError(t, sk.SetNewValidatorByPowerIndex(ctx, val))
	}

	delAcc1 := types.GetRandomBech32Addr()
	delAcc2 := types.GetRandomBech32Addr()
	delAcc3 := types.GetRandomBech32Addr()

	require.NoError(t, sk.SetDelegation(ctx, stakingtypes.NewDelegation(acc1.String(), valAddrs[0].String(), cosmos.NewDec(100))))
	require.NoError(t, sk.SetDelegation(ctx, stakingtypes.NewDelegation(acc2.String(), valAddrs[1].String(), cosmos.NewDec(200))))
	require.NoError(t, sk.SetDelegation(ctx, stakingtypes.NewDelegation(acc3.String(), valAddrs[2].String(), cosmos.NewDec(500))))

	del1 := stakingtypes.NewDelegation(delAcc1.String(), valAddrs[0].String(), cosmos.NewDec(10))
	del2 := stakingtypes.NewDelegation(delAcc2.String(), valAddrs[1].String(), cosmos.NewDec(20))
	del3 := stakingtypes.NewDelegation(delAcc3.String(), valAddrs[2].String(), cosmos.NewDec(20))
	require.NoError(t, sk.SetDelegation(ctx, del1))
	require.NoError(t, sk.SetDelegation(ctx, del2))
	require.NoError(t, sk.SetDelegation(ctx, del3))

	// Mint initial funds to the reserve
	require.NoError(t, k.MintToModule(ctx, types.ModuleName, getCoin(common.Tokens(200000))))

	ctx = ctx.WithBlockHeight(mgr.FetchConfig(ctx, configs.ValidatorPayoutCycle))

	// Create VoteInfo for each validator
	votes := make([]abci.VoteInfo, len(vals))
	for i, val := range vals {
		consAddr, err := val.GetConsAddr()
		require.NoError(t, err)
		votes[i] = abci.VoteInfo{
			Validator: abci.Validator{
				Address: consAddr,
				Power:   val.Tokens.Int64(),
			},
			BlockIdFlag: 2,
		}
	}
	balanceDistribution, err := mgr.keeper.MintAndDistributeTokens(ctx, blockReward)
	if err != nil {
		ctx.Logger().Error("unable to mint and distribute tokens", "error", err)
	}
	devAccountAddress, err := sdk.AccAddressFromBech32(types.FoundationDevAccount)
	require.NoError(t, err)

	grantAccountAddress, err := sdk.AccAddressFromBech32(types.FoundationGrantsAccount)
	require.NoError(t, err)

	communityAccountAddress, err := sdk.AccAddressFromBech32(types.FoundationCommunityAccount)
	require.NoError(t, err)

	devAccountBal := k.GetBalance(ctx, devAccountAddress).AmountOf(configs.Denom)
	require.Equal(t, devAccountBal, sdkmath.NewInt(400000))

	grantAccountBal := k.GetBalance(ctx, grantAccountAddress).AmountOf(configs.Denom)
	require.Equal(t, grantAccountBal, sdkmath.NewInt(400000))

	communityAccountBal := k.GetBalance(ctx, communityAccountAddress).AmountOf(configs.Denom)
	require.Equal(t, communityAccountBal, sdkmath.NewInt(200000))

	moduleBalance := k.GetBalanceOfModule(ctx, types.ModuleName, configs.Denom)
	require.Equal(t, moduleBalance.Int64(), int64(20000001000000))

	require.NoError(t, mgr.ValidatorPayout(ctx, votes, balanceDistribution))

	totalBal := cosmos.ZeroInt()

	// Check balances of validators 7
	checkBalance(ctx, t, k, acc1, configs.Denom, 117529, &totalBal)
	checkBalance(ctx, t, k, acc2, configs.Denom, 234824, &totalBal)
	checkBalance(ctx, t, k, acc3, configs.Denom, 585294, &totalBal)

	// Check balances of delegates
	checkBalance(ctx, t, k, delAcc1, configs.Denom, 11753, &totalBal)
	checkBalance(ctx, t, k, delAcc2, configs.Denom, 23482, &totalBal)
	checkBalance(ctx, t, k, delAcc3, configs.Denom, 23411, &totalBal)

	require.Equal(t, totalBal.ToLegacyDec(), sdkmath.LegacyNewDec(996293))

	moduleBalance = k.GetBalanceOfModule(ctx, types.ModuleName, configs.Denom)
	require.Equal(t, moduleBalance.Int64(), int64(20000000000000))
}
func checkBalance(ctx cosmos.Context, t *testing.T, k Keeper, acc cosmos.AccAddress, denom string, expectedAmt int64, total *sdkmath.Int) {
	bal := k.GetBalance(ctx, acc)
	*total = total.Add(bal.AmountOf(denom))
	require.Equal(t, bal.AmountOf(denom).Int64(), expectedAmt)
}
