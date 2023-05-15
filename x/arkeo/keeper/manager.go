package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"

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

func (mgr *Manager) BeginBlock(ctx cosmos.Context, req abci.RequestBeginBlock) error {
	// if local version is behind the consensus version, panic and don't try to
	// create a new block
	ver := mgr.keeper.GetVersion(ctx)
	if ver > configs.SWVersion {
		panic(
			fmt.Sprintf("Unsupported Version: update your binary (your version: %d, network consensus version: %d)",
				configs.SWVersion,
				ver,
			),
		)
	}

	if err := mgr.ValidatorPayout(ctx, req.LastCommitInfo.GetVotes()); err != nil {
		ctx.Logger().Error("unable to settle contracts", "error", err)
	}
	return nil
}

func (mgr Manager) EndBlock(ctx cosmos.Context) error {
	if err := mgr.ContractEndBlock(ctx); err != nil {
		ctx.Logger().Error("unable to settle contracts", "error", err)
	}

	// invariant checks
	if err := mgr.invariantBondModule(ctx); err != nil {
		panic(err)
	}
	if err := mgr.invariantContractModule(ctx); err != nil {
		panic(err)
	}
	if err := mgr.invariantMaxSupply(ctx); err != nil {
		panic(err)
	}
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
			ctx.Logger().Error("fail to unmarshal provider", "error", err)
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
func (mgr Manager) ValidatorPayout(ctx cosmos.Context, votes []abci.VoteInfo) error {
	valCycle := mgr.FetchConfig(ctx, configs.ValidatorPayoutCycle)
	if valCycle == 0 || ctx.BlockHeight()%valCycle != 0 {
		return nil
	}
	emissionCurve := mgr.FetchConfig(ctx, configs.EmissionCurve)
	blocksPerYear := mgr.FetchConfig(ctx, configs.BlocksPerYear)

	reserveBal := mgr.keeper.GetBalance(ctx, mgr.keeper.GetModuleAccAddress(types.ReserveName))
	for _, bal := range reserveBal {
		reserve := bal.Amount
		blockReward := mgr.calcBlockReward(reserve.Int64(), emissionCurve, (blocksPerYear / valCycle))

		if blockReward.IsZero() {
			continue
		}

		// sum tokens
		total := cosmos.ZeroInt()
		for _, vote := range votes {
			val := mgr.sk.ValidatorByConsAddr(ctx, vote.Validator.Address)
			if val == nil {
				ctx.Logger().Info("unable to find validator", "validator", string(vote.Validator.Address))
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
			if !vote.SignedLastBlock {
				ctx.Logger().Info("validator rewards skipped due to lack of signature", "validator", string(vote.Validator.Address))
				continue
			}

			val := mgr.sk.ValidatorByConsAddr(ctx, vote.Validator.Address)
			if val == nil {
				ctx.Logger().Info("unable to find validator", "validator", string(vote.Validator.Address))
				continue
			}
			if !val.IsBonded() || val.IsJailed() {
				ctx.Logger().Info("validator rewards skipped due to status or jailed", "validator", val.GetOperator().String())
				continue
			}

			acc := cosmos.AccAddress(val.GetOperator())
			valVersion := mgr.keeper.GetStoreVersionForAddress(ctx, acc)
			if valVersion < mgr.keeper.GetVersion(ctx) {
				continue
			}

			totalReward := common.GetSafeShare(val.GetDelegatorShares().RoundInt(), total, blockReward)
			validatorReward := cosmos.ZeroInt()
			rateBasisPts := val.GetCommission().MulInt64(100).RoundInt()

			delegates := mgr.sk.GetValidatorDelegations(ctx, val.GetOperator())
			for _, delegate := range delegates {
				delegateAcc, err := cosmos.AccAddressFromBech32(delegate.DelegatorAddress)
				if err != nil {
					ctx.Logger().Error("unable to fetch delegate address", "delegate", delegate.DelegatorAddress, "error", err)
					continue
				}
				delegateReward := common.GetSafeShare(delegate.GetShares().RoundInt(), val.GetDelegatorShares().RoundInt(), totalReward)
				if acc.String() != delegate.DelegatorAddress {
					valFee := common.GetSafeShare(rateBasisPts, cosmos.NewInt(configs.MaxBasisPoints), delegateReward)
					delegateReward = delegateReward.Sub(valFee)
					validatorReward = validatorReward.Add(valFee)
				}
				if err := mgr.keeper.SendFromModuleToAccount(ctx, types.ReserveName, delegateAcc, cosmos.NewCoins(cosmos.NewCoin(bal.Denom, delegateReward))); err != nil {
					ctx.Logger().Error("unable to pay rewards to delegate", "delegate", delegate.DelegatorAddress, "error", err)
				}
				ctx.Logger().Info("delegate rewarded", "delegate", delegateAcc.String(), "amount", delegateReward)
			}

			if !validatorReward.IsZero() {
				if err := mgr.keeper.SendFromModuleToAccount(ctx, types.ReserveName, acc, cosmos.NewCoins(cosmos.NewCoin(bal.Denom, validatorReward))); err != nil {
					ctx.Logger().Error("unable to pay rewards to validator", "validator", val.GetOperator().String(), "error", err)
					continue
				}
				ctx.Logger().Info("validator additional rewards", "validator", acc.String(), "amount", validatorReward)
			}

			if err := mgr.EmitValidatorPayoutEvent(ctx, acc, validatorReward); err != nil {
				ctx.Logger().Error("unable to emit validator payout event", "validator", acc.String(), "error", err)
			}
		}
	}

	return nil
}

func (mgr Manager) calcBlockReward(totalReserve, emissionCurve, blocksPerYear int64) cosmos.Int {
	// Block Rewards will take the latest reserve, divide it by the emission
	// curve factor, then divide by blocks per year
	if emissionCurve == 0 || blocksPerYear == 0 {
		return cosmos.ZeroInt()
	}
	trD := cosmos.NewDec(totalReserve)
	ecD := cosmos.NewDec(emissionCurve)
	bpyD := cosmos.NewDec(blocksPerYear)
	return trD.Quo(ecD).Quo(bpyD).RoundInt()
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
		if err := mgr.keeper.SendFromModuleToModule(ctx, types.ContractName, types.ReserveName, cosmos.NewCoins(cosmos.NewCoin(contract.Rate.Denom, valIncome))); err != nil {
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
		debt = contract.Rate.Amount.MulRaw(ctx.BlockHeight() - contract.Height).Sub(contract.Paid)
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
