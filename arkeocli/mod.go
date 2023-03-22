package arkeocli

import (
	"strconv"
	"strings"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func newModProviderCmd() *cobra.Command {
	modProviderCmd := &cobra.Command{
		Use:   "mod-provider",
		Short: "mod provider",
		Args:  cobra.ExactArgs(0),
		RunE:  runModProviderCmd,
	}

	flags.AddTxFlagsToCmd(modProviderCmd)
	modProviderCmd.Flags().String("provider-pubkey", "", "provider pubkey")
	modProviderCmd.Flags().String("service", "", "provider service name")
	modProviderCmd.Flags().String("status", "", "provider status (online or offline)")
	modProviderCmd.Flags().String("meta-uri", "", "public endpoint where metadata can be found")
	modProviderCmd.Flags().Uint64("meta-nonce", 0, "increment with each metadata change")
	modProviderCmd.Flags().Uint64("min-duration", 0, "minimum contract duration (in blocks)")
	modProviderCmd.Flags().Uint64("max-duration", 0, "maximum contract duration (in blocks)")
	modProviderCmd.Flags().Uint64("settlement-duration", 0, "settlement duration (in blocks)")
	modProviderCmd.Flags().Uint64("subscription-rate", 0, "rate for subscription contracts")
	modProviderCmd.Flags().Uint64("pay-as-you-go-rate", 0, "rate for pay-as-you-go contracts")
	return modProviderCmd
}

func runModProviderCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	key, err := ensureKeys(cmd)
	if err != nil {
		return
	}
	addr, err := key.GetAddress()
	if err != nil {
		return
	}
	clientCtx = clientCtx.WithFromName(key.Name).WithFromAddress(addr)
	if err = client.SetCmdClientContext(cmd, clientCtx); err != nil {
		return
	}

	argPubkey, _ := cmd.Flags().GetString("provider-pubkey")
	if argPubkey == "" {
		argPubkey, err = toPubkey(cmd, addr)
		if err != nil {
			return
		}
	}

	argService, _ := cmd.Flags().GetString("service")
	if argService == "" {
		argService, err = promptForArg(cmd, "Specify service (e.g. gaia-mainnet-rpc-archive, btc-mainnet-fullnode, etc): ")
		if err != nil {
			return err
		}
	}

	argStatus, _ := cmd.Flags().GetString("status")
	if argStatus == "" {
		argStatus, err = promptForArg(cmd, "Specify status (one of online or offline): ")
		if err != nil {
			return err
		}
	}

	argMetaURI, _ := cmd.Flags().GetString("meta-uri")
	if argMetaURI == "" {
		argMetaURI, err = promptForArg(cmd, "Specify public endpoint where metadata can be found: ")
		if err != nil {
			return err
		}
	}

	argMetaNonce, _ := cmd.Flags().GetUint64("meta-nonce")
	if argMetaNonce == 0 {
		nonce, err := promptForArg(cmd, "Increment nonce to signal provider metadata changed: ")
		if err != nil {
			return err
		}
		argMetaNonce, err = strconv.ParseUint(nonce, 10, 64)
		if err != nil {
			return err
		}
	}

	argMinDuration, _ := cmd.Flags().GetUint64("min-duration")
	if argMinDuration == 0 {
		duration, err := promptForArg(cmd, "Specify minimum contract duration (in blocks): ")
		if err != nil {
			return err
		}
		argMinDuration, err = strconv.ParseUint(duration, 10, 64)
		if err != nil {
			return err
		}
	}

	argMaxDuration, _ := cmd.Flags().GetUint64("max-duration")
	if argMaxDuration == 0 {
		duration, err := promptForArg(cmd, "Specify maximum contract duration (in blocks): ")
		if err != nil {
			return err
		}
		argMaxDuration, err = strconv.ParseUint(duration, 10, 64)
		if err != nil {
			return err
		}
	}

	argSettlementDuration, _ := cmd.Flags().GetUint64("settlement-duration")
	if argSettlementDuration == 0 {
		duration, err := promptForArg(cmd, "Specify settlement duration (in blocks): ")
		if err != nil {
			return err
		}
		argSettlementDuration, err = strconv.ParseUint(duration, 10, 64)
		if err != nil {
			return err
		}
	}

	argSubscriptionRate, _ := cmd.Flags().GetUint64("subscription-rate")
	if argSubscriptionRate == 0 {
		subscriptiomRate, err := promptForArg(cmd, "Specify rate for subscription contracts: ")
		if err != nil {
			return err
		}
		argSubscriptionRate, err = strconv.ParseUint(subscriptiomRate, 10, 64)
		if err != nil {
			return err
		}
	}

	argPayAsYouGoRate, _ := cmd.Flags().GetUint64("pay-as-you-go-rate")
	if argPayAsYouGoRate == 0 {
		payAsYouGoRate, err := promptForArg(cmd, "Specify rate for pay-as-you-go contracts: ")
		if err != nil {
			return err
		}
		argPayAsYouGoRate, err = strconv.ParseUint(payAsYouGoRate, 10, 64)
		if err != nil {
			return err
		}
	}

	pubkey, err := common.NewPubKey(argPubkey)
	if err != nil {
		return err
	}

	status := types.ProviderStatus(types.ProviderStatus_value[strings.ToUpper(argStatus)])

	msg := types.NewMsgModProvider(
		clientCtx.GetFromAddress().String(),
		pubkey,
		argService,
		argMetaURI,
		argMetaNonce,
		status,
		int64(argMinDuration),
		int64(argMaxDuration),
		int64(argSubscriptionRate),
		int64(argPayAsYouGoRate),
		int64(argSettlementDuration),
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
