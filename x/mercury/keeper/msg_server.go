package keeper

import (
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"
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

// convert int64s into coins asset
func getCoins(vals ...int64) cosmos.Coins {
	coins := make(cosmos.Coins, len(vals))
	for i, val := range vals {
		coins[i] = getCoin(val)
	}
	return coins
}
*/

// convert int64 into coin asset
func getCoin(val int64) cosmos.Coin {
	return cosmos.NewCoin(configs.Denom, cosmos.NewInt(val))
}

func tokens(i int64) int64 {
	return i * (10 * cosmos.DefaultCoinDecimals)
}
