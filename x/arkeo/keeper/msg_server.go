package keeper

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type msgServer struct {
	Keeper
	mgr Manager
}

func newMsgServer(keeper Keeper, sk stakingkeeper.Keeper) *msgServer {
	return &msgServer{
		Keeper: keeper,
		mgr:    NewManager(keeper, sk),
	}
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, sk stakingkeeper.Keeper) types.MsgServer {
	return newMsgServer(keeper, sk)
}

var _ types.MsgServer = msgServer{}

func (k msgServer) FetchConfig(ctx cosmos.Context, name configs.ConfigName) int64 {
	// TODO: use ctx to fetch config overrides from the chain state
	return k.mgr.Configs(ctx).GetInt64Value(name)
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
