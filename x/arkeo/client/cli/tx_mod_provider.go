package cli

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdModProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mod-provider [pubkey] [service] [metatadata-uri] [metadata-nonce] [status] [min-contract-duration] [max-contract-duration] [pay-per-block-rate] [pay-per-call-rate] [settlement-duration]",
		Short: "Broadcast message modProvider",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}
			argService := args[1]

			argMetatadataURI := args[2]
			argMetadataNonce, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}
			argStatus, err := cast.ToInt32E(args[4])
			if err != nil {
				return err
			}
			argMinContractDuration, err := cast.ToInt64E(args[5])
			if err != nil {
				return err
			}
			argMaxContractDuration, err := cast.ToInt64E(args[6])
			if err != nil {
				return err
			}
			argPayPerBlockRate, err := cast.ToInt64E(args[7])
			if err != nil {
				return err
			}
			argPayPeryCallRate, err := cast.ToInt64E(args[8])
			if err != nil {
				return err
			}

			argSettlementDuration, err := cast.ToInt64E(args[9])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rates := []*types.ContractRate{
				{
					UserType:  types.UserType_SINGLE_USER,
					MeterType: types.MeterType_PAY_PER_BLOCK,
					Rate:      argPayPerBlockRate,
				},
				{
					UserType:  types.UserType_SINGLE_USER,
					MeterType: types.MeterType_PAY_PER_CALL,
					Rate:      argPayPeryCallRate,
				},
			}

			msg := types.NewMsgModProvider(
				clientCtx.GetFromAddress().String(),
				pubkey,
				argService,
				argMetatadataURI,
				argMetadataNonce,
				types.ProviderStatus(argStatus),
				argMinContractDuration,
				argMaxContractDuration,
				argSettlementDuration,
				rates,
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
