package cli

import (
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdOpenContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-contract [provider_pubkey] [service] [client_pubkey] [user-type] [meter-type] [deposit] [duration] [rate] [settlement-duration] [delegation-optional]",
		Short: "Broadcast message openContract",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argService := args[1]
			argClient := args[2]
			argDeposit := args[5]

			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}

			cl, err := common.NewPubKey(argClient)
			if err != nil {
				return err
			}

			argUserType, err := cast.ToInt32E(args[3])
			if err != nil {
				return err
			}

			argMeterType, err := cast.ToInt32E(args[4])
			if err != nil {
				return err
			}

			argDuration, err := cast.ToInt64E(args[5])
			if err != nil {
				return err
			}

			argRate, err := cosmos.ParseCoin(args[6])
			if err != nil {
				return err
			}

			argSettlementDuration, err := cast.ToInt64E(args[7])
			if err != nil {
				return err
			}

			deposit, ok := cosmos.NewIntFromString(argDeposit)
			if !ok {
				return fmt.Errorf("bad deposit amount: %s", argDeposit)
			}

			delegate := common.EmptyPubKey
			if len(args) > 8 {
				delegate, err = common.NewPubKey(args[8])
				if err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgOpenContract(
				clientCtx.GetFromAddress(),
				pubkey,
				argService,
				cl,
				delegate,
				types.UserType(argUserType),
				types.MeterType(argMeterType),
				argDuration,
				argSettlementDuration,
				argRate,
				deposit,
				types.Restrictions{}, // TODO: add restrictions
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
