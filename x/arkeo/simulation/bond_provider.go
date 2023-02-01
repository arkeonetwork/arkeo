package simulation

import (
	"math/rand"

	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/keeper"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgBondProvider(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgBondProvider{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the BondProvider simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "BondProvider simulation not implemented"), nil, nil
	}
}
