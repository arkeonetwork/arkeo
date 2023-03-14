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
		Use:   "open-contract [provider_pubkey] [service] [client_pubkey] [c-type] [deposit] [duration] [rate] [settlement-duration] [delegation-optional]",
		Short: "Broadcast message openContract",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argChain := args[1]
			argClient := args[2]
			argDeposit := args[4]

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

			argRate, err := cast.ToInt64E(args[6])
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
				clientCtx.GetFromAddress().String(),
				pubkey,
				argChain,
				cl,
				delegate,
				types.ContractType(argContractType),
				argDuration,
				argSettlementDuration,
				argRate,
				deposit,
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
