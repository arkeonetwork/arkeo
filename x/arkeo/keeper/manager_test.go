package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestValidatorPayout(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	pks := simapp.CreateTestPubKeys(3)
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

	valAddrs := simapp.ConvertAddrsToValAddrs([]cosmos.AccAddress{acc1, acc2, acc3})

	val1, err := stakingtypes.NewValidator(valAddrs[0], pks[0], stakingtypes.Description{})
	require.NoError(t, err)
	val1.Tokens = cosmos.NewInt(100)
	val1.DelegatorShares = cosmos.NewDec(100 + 10 + 20)
	val1.Status = stakingtypes.Bonded
	val1.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(1, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	val2, err := stakingtypes.NewValidator(valAddrs[1], pks[1], stakingtypes.Description{})
	require.NoError(t, err)
	val2.Tokens = cosmos.NewInt(200)
	val2.DelegatorShares = cosmos.NewDec(200 + 20)
	val2.Status = stakingtypes.Bonded
	val2.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(2, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	val3, err := stakingtypes.NewValidator(valAddrs[2], pks[2], stakingtypes.Description{})
	require.NoError(t, err)
	val3.Tokens = cosmos.NewInt(500)
	val3.DelegatorShares = cosmos.NewDec(500)
	val3.Status = stakingtypes.Bonded
	val3.Commission = stakingtypes.NewCommission(cosmos.NewDecWithPrec(5, 1), cosmos.ZeroDec(), cosmos.ZeroDec())

	vals := []stakingtypes.Validator{val1, val2, val3}
	for _, val := range vals {
		sk.SetValidator(ctx, val)
		require.NoError(t, sk.SetValidatorByConsAddr(ctx, val))
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

	require.NoError(t, k.MintToModule(ctx, types.ModuleName, getCoin(common.Tokens(50000))))
	require.NoError(t, k.SendFromModuleToModule(ctx, types.ModuleName, types.ReserveName, getCoins(common.Tokens(50000))))

	mgr := NewManager(k, sk)
	ctx = ctx.WithBlockHeight(mgr.FetchConfig(ctx, configs.ValidatorPayoutCycle))

	blockReward := int64(237792)
	require.NoError(t, mgr.ValidatorEndBlock(ctx))
	require.Equal(t, k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), 5000000000000-blockReward)

	// check validator balances
	totalBal := cosmos.ZeroInt()
	bal := k.GetBalance(ctx, acc1)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(27984))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, acc2)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(55962))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, acc3)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(139878))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	// check delegate balances
	bal = k.GetBalance(ctx, delAcc1)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(2795))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, delAcc2)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(5589))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	bal = k.GetBalance(ctx, delAcc3)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(5584))
	totalBal = totalBal.Add(bal.AmountOf(configs.Denom))

	// ensure block reward is equal to total rewarded to validators and delegates
	require.Equal(t, blockReward, totalBal.Int64())
}

func TestContractEndBlock(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)
	mgr := NewManager(k, sk)

	// create a provider for 2 services
	providerPubKey := types.GetRandomPubKey()
	provider := types.NewProvider(providerPubKey, common.BTCService)
	provider.Bond = cosmos.NewInt(20000000000)
	provider.LastUpdate = ctx.BlockHeight()
	require.NoError(t, k.SetProvider(ctx, provider))
	provider.Service = common.ETHService
	require.NoError(t, k.SetProvider(ctx, provider))

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Service:             common.BTCService.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		Rates: []*types.ContractRate{
			{
				UserType:  types.UserType_SINGLE_USER,
				MeterType: types.MeterType_PAY_PER_BLOCK,
				Rate:      15,
			},
			{
				UserType:  types.UserType_SINGLE_USER,
				MeterType: types.MeterType_PAY_PER_CALL,
				Rate:      15,
			},
		},
	}
	err := s.ModProviderHandle(ctx, &modProviderMsg)
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
		Provider:  providerPubKey,
		Service:   common.BTCService.String(),
		Creator:   user1Address.String(),
		Client:    user1PubKey,
		MeterType: types.MeterType_PAY_PER_CALL,
		UserType:  types.UserType_SINGLE_USER,
		Duration:  100,
		Rate:      15,
		Deposit:   cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// have user1 open a contract for a delegate.
	delegatePubKey := types.GetRandomPubKey()
	msg.Delegate = delegatePubKey
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// create user2 to open a contract against the provider.
	user2PubKey := types.GetRandomPubKey()
	user2Address, err := user2PubKey.GetMyAddress()
	require.NoError(t, err)

	require.NoError(t, k.MintAndSendToAccount(ctx, user2Address, getCoin(common.Tokens(20))))
	msg.Delegate = common.EmptyPubKey
	msg.Client = user2PubKey
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

	// create a provider for 2 services
	providerPubKey := types.GetRandomPubKey()
	provider := types.NewProvider(providerPubKey, common.BTCService)
	provider.Bond = cosmos.NewInt(20000000000)
	provider.LastUpdate = ctx.BlockHeight()
	provider.SettlementDuration = 10
	require.NoError(t, k.SetProvider(ctx, provider))
	provider.Service = common.ETHService
	require.NoError(t, k.SetProvider(ctx, provider))

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Service:             common.BTCService.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		Rates: []*types.ContractRate{
			{
				UserType:  types.UserType_SINGLE_USER,
				MeterType: types.MeterType_PAY_PER_BLOCK,
				Rate:      15,
			},
			{
				UserType:  types.UserType_SINGLE_USER,
				MeterType: types.MeterType_PAY_PER_CALL,
				Rate:      15,
			},
		},
		SettlementDuration: 10,
	}

	err := s.ModProviderHandle(ctx, &modProviderMsg)
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
		Provider:           providerPubKey,
		Service:            common.BTCService.String(),
		Creator:            user1Address.String(),
		Client:             user1PubKey,
		UserType:           types.UserType_SINGLE_USER,
		MeterType:          types.MeterType_PAY_PER_CALL,
		Duration:           100,
		Rate:               15,
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
