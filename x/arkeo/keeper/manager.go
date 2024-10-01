package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmptm "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type Manager struct {
	keeper Keeper
	sk     stakingkeeper.Keeper
}

func NewManager(k Keeper, sk stakingkeeper.Keeper) Manager {
	return Manager{
		keeper: k,
		sk:     sk,
	}
}

func (mgr *Manager) BeginBlock(ctx cosmos.Context) error {
	// if local version is behind the consensus version, panic and don't try to
	// create a new block

	params := mgr.keeper.GetParams(ctx)
	ver := mgr.keeper.GetComputedVersion(ctx)
	swVersion, err := configs.GetSWVersion()
	if err != nil {
		return err
	}
	if ver > swVersion {
		panic(
			fmt.Sprintf("Unsupported Version: update your binary (your version: %d, network consensus version: %d)",
				swVersion,
				ver,
			),
		)
	}
	mgr.keeper.SetVersion(ctx, ver) // update stored version

	// Get the circulating supply after calculating inflation
	circSupply, err := mgr.circulatingSupplyAfterInflationCalc(ctx)
	if err != nil {
		mgr.keeper.Logger().Error("unable to get supply with inflation calculation", "error", err)
		return err
	}

	err = mgr.keeper.MoveTokensFromDistributionToFoundationPoolAccount(ctx)
	if err != nil {
		mgr.keeper.Logger().Error("unable to send tokens from distribution to pool account", "error", err)
	}

	validatorPayoutCycle := sdkmath.LegacyNewDec(mgr.FetchConfig(ctx, configs.ValidatorPayoutCycle))

	emissionCurve := sdkmath.LegacyNewDec(int64(params.EmissionCurve)) // Emission curve factor
	blocksPerYear := sdkmath.LegacyNewDec(int64(params.BlockPerYear))

	// Distribute Minted To Pools
	balanceDistribution, err := mgr.keeper.MintAndDistributeTokens(ctx, circSupply)
	if err != nil {
		mgr.keeper.Logger().Error("unable to mint and distribute tokens", "error", err)
	}

	mgr.keeper.Logger().Info(fmt.Sprintf("Circulating Supply After Funding Foundational Account  %s", balanceDistribution))

	// Calculate Block Rewards
	blockReward := mgr.calcBlockReward(ctx, balanceDistribution.Amount, emissionCurve, blocksPerYear, validatorPayoutCycle)
	mgr.keeper.Logger().Info(fmt.Sprintf("Block Reward for block number %d, %v", ctx.BlockHeight(), blockReward))

	var votes = []abci.VoteInfo{}
	for i := 0; i < ctx.CometInfo().GetLastCommit().Votes().Len(); i++ {
		vote := ctx.CometInfo().GetLastCommit().Votes().Get(i)
		abciVote := abci.VoteInfo{
			Validator: abci.Validator{
				Address: vote.Validator().Address(),
				Power:   vote.Validator().Power(),
			},
			BlockIdFlag: cmptm.BlockIDFlag(vote.GetBlockIDFlag()),
		}
		votes = append(votes, abciVote)
	}

	if err := mgr.ValidatorPayout(ctx, votes, blockReward); err != nil {
		mgr.keeper.Logger().Error("unable to settle contracts", "error", err)
	}
	return nil
}

func (mgr Manager) EndBlock(ctx cosmos.Context) error {
	if err := mgr.ContractEndBlock(ctx); err != nil {
		mgr.keeper.Logger().Error("unable to settle contracts", "error", err)
	}

	// invariant checks
	if err := mgr.invariantBondModule(ctx); err != nil {
		panic(err)
	}
	if err := mgr.invariantContractModule(ctx); err != nil {
		panic(err)
	}
	// if err := mgr.invariantMaxSupply(ctx); err != nil {
	// 	panic(err)
	// }
	return nil
}

func (mgr Manager) Configs(ctx cosmos.Context) configs.ConfigValues {
	return configs.GetConfigValues(mgr.keeper.GetVersion(ctx))
}

// test that the bond module has enough bond in it
func (mgr Manager) invariantBondModule(ctx cosmos.Context) error {
	balance := mgr.keeper.GetBalanceOfModule(ctx, types.ProviderName, configs.Denom)

	sum := cosmos.ZeroInt()
	iter := mgr.keeper.GetProviderIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var provider types.Provider
		if err := mgr.keeper.Cdc().Unmarshal(iter.Value(), &provider); err != nil {
			mgr.keeper.Logger().Error("fail to unmarshal provider", "error", err)
			continue
		}
		sum = sum.Add(provider.Bond)
	}

	if sum.GT(balance) {
		// TODO: instead of returning an error and causing a panic, pause the bond provider handler and allow the chain to continue to function
		return errors.Wrapf(types.ErrInvariantBondModule, "bond module does not have enough token in it to back the bond records (%s/%s)", sum.String(), balance.String())
	}

	return nil
}

// test that the contract module has enough bond in it
func (mgr Manager) invariantContractModule(ctx cosmos.Context) error {
	sums := cosmos.NewCoins()
	iter := mgr.keeper.GetContractIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var contract types.Contract
		if err := mgr.keeper.Cdc().Unmarshal(iter.Value(), &contract); err != nil {
			ctx.Logger().Error("fail to unmarshal contract", "error", err)
			continue
		}
		if contract.IsSettled(ctx.BlockHeight()) {
			continue
		}
		sums = sums.Add(cosmos.NewCoin(contract.Rate.Denom, contract.Deposit.Sub(contract.Paid)))
	}

	for _, sum := range sums {
		if sum.Amount.IsZero() {
			continue
		}
		balance := mgr.keeper.GetBalanceOfModule(ctx, types.ContractName, sum.Denom)
		if sum.Amount.GT(balance) {
			return errors.Wrapf(types.ErrInvariantContractModule, "contract module does not have enough token (%s) in it to back the bond records (%s/%s)", sum.Denom, sum.Amount.String(), balance.String())
		}
	}

	return nil
}

// test that the total supply does not surpass max supply
func (mgr Manager) invariantMaxSupply(ctx cosmos.Context) error {
	supply := mgr.keeper.GetSupply(ctx, configs.Denom)

	max := cosmos.NewInt(mgr.FetchConfig(ctx, configs.MaxSupply))
	if supply.Amount.GT(max) {
		return errors.Wrapf(types.ErrInvariantMaxSupply, "supply has surpass the max set (%s/%s)", supply.Amount.String(), max.String())
	}

	return nil
}

func (mgr Manager) ContractEndBlock(ctx cosmos.Context) error {
	set, err := mgr.keeper.GetContractExpirationSet(ctx, ctx.BlockHeight())
	if err != nil {
		return err
	}

	if set.ContractSet == nil || len(set.ContractSet.ContractIds) == 0 {
		return nil
	}

	for _, contractId := range set.ContractSet.ContractIds {
		contract, err := mgr.keeper.GetContract(ctx, contractId)
		if err != nil {
			ctx.Logger().Error("unable to fetch contract", "id", contractId, "error", err)
			continue
		}
		if contract.Client.IsEmpty() {
			continue
		}

		_, err = mgr.SettleContract(ctx, contract, 0, true)
		if err != nil {
			ctx.Logger().Error("unable to settle contract", "id", contractId, "error", err)
			continue
		}
	}

	return nil
}

// This function pays out rewards to validators.
// TODO: the method of accomplishing this is admittedly quite inefficient. The
// better approach would be to track live allocation via assigning "units" to
// validators when they bond and unbond. The math for this is as follows
// U = total bond units
// T = tokens bonded
// t = new tokens being bonded
// units = U / (T / t)
// Since the development goal at the moment is to get this chain up and
// running, we can save this optimization for another day.
func (mgr Manager) ValidatorPayout(ctx cosmos.Context, votes []abci.VoteInfo, blockReward sdk.DecCoin) error {
	if blockReward.IsZero() {
		return nil
	}

	// sum tokens
	total := cosmos.ZeroInt()
	for _, vote := range votes {
		val, err := mgr.sk.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if err != nil {
			mgr.keeper.Logger().Info("unable to find validator", "validator", string(vote.Validator.Address))
			continue
		}
		if !val.IsBonded() || val.IsJailed() {
			continue
		}
		total = total.Add(val.GetDelegatorShares().RoundInt())
	}
	if total.IsZero() {
		return nil
	}

	for _, vote := range votes {
		if vote.BlockIdFlag.String() == "BLOCK_ID_FLAG_ABSENT" || vote.BlockIdFlag.String() == "BLOCK_ID_FLAG_UNKNOWN" {
			mgr.keeper.Logger().Info(fmt.Sprintf("validator rewards skipped due to lack of signature: %s, validator : %s ", vote.BlockIdFlag.String(), string(vote.Validator.GetAddress())))
			continue
		}

		val, err := mgr.sk.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if err != nil {
			mgr.keeper.Logger().Info("unable to find validator", "validator", string(vote.Validator.Address))
			continue
		}
		if !val.IsBonded() || val.IsJailed() {
			mgr.keeper.Logger().Info("validator rewards skipped due to status or jailed", "validator", val.GetOperator())
			continue
		}

		valBz, err := mgr.sk.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			panic(err)
		}

		valVersion := mgr.keeper.GetVersionForAddress(ctx, valBz)
		curVersion := mgr.keeper.GetVersion(ctx)
		if valVersion < curVersion {
			continue
		}
		acc := cosmos.AccAddress(val.GetOperator())

		totalReward := common.GetSafeShare(val.GetDelegatorShares().RoundInt(), total, blockReward.Amount.RoundInt())
		validatorReward := cosmos.ZeroInt()
		rateBasisPts := val.GetCommission().MulInt64(100).RoundInt()

		delegates, err := mgr.sk.GetValidatorDelegations(ctx, valBz)
		if err != nil {
			panic(err)
		}

		for _, delegate := range delegates {
			delegateAcc, err := cosmos.AccAddressFromBech32(delegate.DelegatorAddress)
			if err != nil {
				mgr.keeper.Logger().Error("unable to fetch delegate address", "delegate", delegate.DelegatorAddress, "error", err)
				continue
			}
			delegateReward := common.GetSafeShare(delegate.GetShares().RoundInt(), val.GetDelegatorShares().RoundInt(), totalReward)
			if acc.String() != delegate.DelegatorAddress {
				valFee := common.GetSafeShare(rateBasisPts, cosmos.NewInt(configs.MaxBasisPoints), delegateReward)
				delegateReward = delegateReward.Sub(valFee)
				validatorReward = validatorReward.Add(valFee)
			}

			if err := mgr.keeper.MintAndSendToAccount(ctx, delegateAcc, cosmos.NewCoin(blockReward.Denom, delegateReward)); err != nil {
				mgr.keeper.Logger().Error("unable to pay rewards to delegate", "delegate", delegate.DelegatorAddress, "error", err)
				continue
			}
			mgr.keeper.Logger().Info("delegate rewarded", "delegate", delegateAcc.String(), "amount", delegateReward)
		}

		if !validatorReward.IsZero() {
			if err := mgr.keeper.AllocateTokensToValidator(ctx, val, sdk.NewDecCoins(sdk.NewDecCoin(blockReward.Denom, validatorReward))); err != nil {
				mgr.keeper.Logger().Error("unable to pay rewards to validator", "validator", val.GetOperator(), "error", err)
				continue
			}
			mgr.keeper.Logger().Info("validator additional rewards", "validator", acc.String(), "amount", validatorReward)
		}

		if err := mgr.EmitValidatorPayoutEvent(ctx, acc, validatorReward); err != nil {
			mgr.keeper.Logger().Error("unable to emit validator payout event", "validator", acc.String(), "error", err)
		}
	}

	return nil
}

func (mgr Manager) calcBlockReward(ctx cosmos.Context, totalReserve, emissionCurve, blocksPerYear, validatorPayoutCycle sdkmath.LegacyDec) sdk.DecCoin {
	// Block Rewards will take the latest reserve, divide it by the emission
	// curve factor, then divide by blocks per year
	if emissionCurve.IsZero() || blocksPerYear.IsZero() {
		mgr.keeper.Logger().Info("block and emission-curve cannot be zero")
		return sdk.NewDecCoin(configs.Denom, sdkmath.NewInt(0))
	}

	if validatorPayoutCycle.IsZero() || ctx.BlockHeight()%validatorPayoutCycle.RoundInt64() != 0 {
		mgr.keeper.Logger().Info("validator payout cycle cannot be zero")
		return sdk.NewDecCoin(configs.Denom, sdkmath.NewInt(0))
	}

	bpyD := blocksPerYear.Quo(validatorPayoutCycle)

	blockReward := totalReserve.Quo(emissionCurve).Quo(bpyD).RoundInt()

	return sdk.NewDecCoin(configs.Denom, blockReward)
}

func (mgr Manager) FetchConfig(ctx cosmos.Context, name configs.ConfigName) int64 {
	// TODO: create a handler for admins to be able to change configs on the
	// fly and check them here before returning
	return mgr.Configs(ctx).GetInt64Value(name)
}

// any owed debt is paid to data provider
func (mgr Manager) SettleContract(ctx cosmos.Context, contract types.Contract, nonce int64, isFinal bool) (types.Contract, error) {
	if nonce > contract.Nonce {
		contract.Nonce = nonce
	}
	totalDebt, err := mgr.contractDebt(ctx, contract)
	valIncome := common.GetSafeShare(cosmos.NewInt(mgr.FetchConfig(ctx, configs.ReserveTax)), cosmos.NewInt(configs.MaxBasisPoints), totalDebt)
	debt := totalDebt.Sub(valIncome)
	if err != nil {
		return contract, err
	}
	if !debt.IsZero() {
		provider, err := contract.Provider.GetMyAddress()
		if err != nil {
			return contract, err
		}
		if err := mgr.keeper.SendFromModuleToAccount(ctx, types.ContractName, provider, cosmos.NewCoins(cosmos.NewCoin(contract.Rate.Denom, debt))); err != nil {
			return contract, err
		}
		if err := mgr.keeper.SendFromModuleToModule(ctx, types.ContractName, types.ModuleName, cosmos.NewCoins(cosmos.NewCoin(contract.Rate.Denom, valIncome))); err != nil {
			return contract, err
		}
	}

	contract.Paid = contract.Paid.Add(totalDebt)
	if isFinal {
		remainder := contract.Deposit.Sub(contract.Paid)
		if !remainder.IsZero() {
			client, err := contract.Client.GetMyAddress()
			if err != nil {
				return contract, err
			}
			if err := mgr.keeper.SendFromModuleToAccount(ctx, types.ContractName, client, cosmos.NewCoins(cosmos.NewCoin(contract.Rate.Denom, remainder))); err != nil {
				return contract, err
			}
			// now that the user has some of their funds refunded, the deposit
			// amount should be updated (to reflect that). This also sets Paid
			// == Deposit, which causes the record to be deleted, conserving
			// space
			contract.Deposit = contract.Paid
		}
		contract.SettlementHeight = ctx.BlockHeight()
		// this contract can now be removed from the users list of contracts
		err = mgr.keeper.RemoveFromUserContractSet(ctx, contract.GetSpender(), contract.Id)
		if err != nil {
			return contract, err
		}
	}

	err = mgr.keeper.SetContract(ctx, contract)
	if err != nil {
		return contract, err
	}

	if err = mgr.EmitContractSettlementEvent(ctx, totalDebt, valIncome, &contract); err != nil {
		return contract, err
	}

	return contract, nil
}

func (mgr Manager) contractDebt(ctx cosmos.Context, contract types.Contract) (cosmos.Int, error) {
	var debt cosmos.Int
	switch contract.Type {
	case types.ContractType_SUBSCRIPTION:
		height := ctx.BlockHeight()
		if height > contract.SettlementPeriodEnd() {
			height = contract.SettlementPeriodEnd()
		}
		debt = contract.Rate.Amount.MulRaw(height - contract.Height).Sub(contract.Paid)
	case types.ContractType_PAY_AS_YOU_GO:
		debt = contract.Rate.Amount.MulRaw(contract.Nonce).Sub(contract.Paid)
	default:
		return cosmos.ZeroInt(), errors.Wrapf(types.ErrInvalidContractType, "%s", contract.Type.String())
	}

	if debt.IsNegative() {
		return cosmos.ZeroInt(), nil
	}

	// sanity check, ensure provider cannot take more than deposited into the contract
	if contract.Paid.Add(debt).GT(contract.Deposit) {
		return contract.Deposit.Sub(contract.Paid), nil
	}

	return debt, nil
}

func (mgr Manager) circulatingSupplyAfterInflationCalc(ctx cosmos.Context) (sdk.DecCoin, error) {
	// Get the circulating supply
	circulatingSupply, err := mgr.keeper.GetCirculatingSupply(ctx, configs.Denom)
	if err != nil {
		mgr.keeper.Logger().Error(fmt.Sprintf("failed to get circulating supply %s", err))
		return sdk.NewDecCoin(configs.Denom, sdkmath.NewInt(0)), err
	}

	// Get the inflation rate
	inflationRate, err := mgr.keeper.GetInflationRate(ctx)
	if err != nil {
		return sdk.NewDecCoin(configs.Denom, sdkmath.NewInt(0)), err
	}
	mgr.keeper.Logger().Info(fmt.Sprintf("inflation rate: %d", inflationRate))

	// Multiply circulating supply by inflation rate to get the newly minted token amount
	newTokenAmountMintedDec := circulatingSupply.Amount.Mul(inflationRate).QuoInt64(100)

	mgr.keeper.Logger().Info(fmt.Sprintf("After Inflation Calculation: %v", newTokenAmountMintedDec))

	return sdk.NewDecCoin(configs.Denom, newTokenAmountMintedDec.RoundInt()), nil
}
