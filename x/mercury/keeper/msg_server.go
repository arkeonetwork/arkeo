package keeper

import (
	"fmt"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
	configs configs.ConfigValues
}

func newMsgServer(keeper Keeper) *msgServer {
	ver := keeper.GetVersion()
	return &msgServer{
		Keeper:  keeper,
		configs: configs.GetConfigValues(ver),
	}
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return newMsgServer(keeper)
}

var _ types.MsgServer = msgServer{}

func (k msgServer) FetchConfig(ctx cosmos.Context, name configs.ConfigName) int64 {
	// TODO: use ctx to fetch config overrides from the chain state
	return k.configs.GetInt64Value(name)
}

// any owed debt is paid to data provider
func (k msgServer) SettleContract(ctx cosmos.Context, contract types.Contract, closed bool) (types.Contract, error) {
	debt, err := k.ContractDebt(ctx, contract)
	fmt.Printf("Debt: %d\n", debt.Int64())
	if err != nil {
		return contract, err
	}
	if !debt.IsZero() {
		provider, err := contract.ProviderPubKey.GetMyAddress()
		if err != nil {
			return contract, err
		}
		if err := k.SendFromModuleToAccount(ctx, types.ContractName, provider, cosmos.NewCoins(cosmos.NewCoin(configs.Denom, debt))); err != nil {
			return contract, err
		}
	}

	contract.Paid = contract.Paid.Add(debt)
	if closed {
		contract.ClosedHeight = ctx.BlockHeight()
	}

	fmt.Printf("Contract: %+v\n", contract)
	err = k.SetContract(ctx, contract)
	if err != nil {
		return contract, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeContractSettlement,
				sdk.NewAttribute("pubkey", contract.ProviderPubKey.String()),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.ClientAddress.String()),
				sdk.NewAttribute("paid", debt.String()),
			),
		},
	)
	return contract, nil
}

func (k msgServer) ContractDebt(ctx cosmos.Context, contract types.Contract) (cosmos.Int, error) {
	var debt cosmos.Int
	switch contract.Type {
	case types.ContractType_Subscription:
		debt = cosmos.NewInt(contract.Rate * (ctx.BlockHeight() - contract.Height)).Sub(contract.Paid)
	case types.ContractType_PayAsYouGo:
		debt = cosmos.NewInt(contract.Rate * contract.Queries).Sub(contract.Paid)
	default:
		return cosmos.ZeroInt(), sdkerrors.Wrapf(types.ErrInvalidContractType, "%s", contract.Type.String())
	}

	if debt.IsNegative() {
		return cosmos.ZeroInt(), nil
	}
	return debt, nil
}

/*
func (k msgServer) getFee(ctx cosmos.Context, names ...configs.ConfigName) int64 {
	var total int64
	for _, name := range names {
		total += k.FetchConfig(ctx, name)
	}
	return total
}

func (k msgServer) hasCoins(ctx cosmos.Context, addr cosmos.AccAddress, names ...configs.ConfigName) error {
	total := k.getFee(ctx, names...)
	coins := getCoins(total)
	if !k.HasCoins(ctx, addr, coins) {
		return sdkerrors.Wrapf(types.ErrInsufficientFunds, "insufficient funds")
	}
	return nil
}
*/

// convert int64s into coins asset
func getCoins(vals ...int64) cosmos.Coins {
	coins := make(cosmos.Coins, len(vals))
	for i, val := range vals {
		coins[i] = getCoin(val)
	}
	return coins
}

// convert int64 into coin asset
func getCoin(val int64) cosmos.Coin {
	return cosmos.NewCoin(configs.Denom, cosmos.NewInt(val))
}
