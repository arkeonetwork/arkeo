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
		Use:   "open-contract [provider_pubkey] [service] [client_pubkey] [c-type] [deposit] [duration] [rate] [queries-per-minute] [settlement-duration] [authorization-optional] [delegation-optional]",
		Short: "Broadcast message openContract",
		Args:  cobra.MinimumNArgs(9),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argService := args[1]
			argClient := args[2]
			argDeposit := args[4]
			deposit, ok := cosmos.NewIntFromString(argDeposit)
			if !ok {
				return fmt.Errorf("bad deposit amount: %s", argDeposit)
			}

			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}

			cl, err := common.NewPubKey(argClient)
			if err != nil {
				return err
			}

			argContractType, err := cast.ToInt32E(args[3])
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

			argQPM, err := cast.ToInt64E(args[7])
			if err != nil {
				return err
			}

			argSettlementDuration, err := cast.ToInt64E(args[8])
			if err != nil {
				return err
			}

			argContractAuth := int32(0)
			if len(args) > 9 {
				argContractAuth, err = cast.ToInt32E(args[9])
				if err != nil {
					return err
				}
			}

			delegate := common.EmptyPubKey
			if len(args) > 10 {
				delegate, err = common.NewPubKey(args[10])
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
				types.ContractType(argContractType),
				argDuration,
				argSettlementDuration,
				argRate,
				deposit,
				types.ContractAuthorization(argContractAuth),
				argQPM,
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
