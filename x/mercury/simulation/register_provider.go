package simulation

import (
	"math/rand"

	"mercury/x/mercury/keeper"
	"mercury/x/mercury/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgRegisterProvider(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRegisterProvider{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the RegisterProvider simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "RegisterProvider simulation not implemented"), nil, nil
	}
}
