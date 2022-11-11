package keeper

import (
	"arkeo/common/cosmos"
	"arkeo/x/arkeo/configs"
	"arkeo/x/arkeo/types"
)

type msgServer struct {
	Keeper
	mgr     Manager
	configs configs.ConfigValues
}

func newMsgServer(keeper Keeper) *msgServer {
	ver := keeper.GetVersion()
	return &msgServer{
		Keeper:  keeper,
		mgr:     NewManager(keeper),
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
