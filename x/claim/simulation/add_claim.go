package simulation

import (
	"math/rand"

	"github.com/arkeonetwork/arkeo/x/claim/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddClaim(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddClaim{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the AddClaim simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "AddClaim simulation not implemented"), nil, nil
	}
}
