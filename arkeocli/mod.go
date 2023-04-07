package arkeocli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var (
	payments  *[]string
	rates     *[]string
	meters    *[]string
	userTypes *[]string
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

	payments = modProviderCmd.Flags().StringArray("payment", []string{}, "adds an accepted payment currency. rate, meter, user-type must be specified after each occurence of this flag")
	rates = modProviderCmd.Flags().StringArray("rate", []string{}, "adds a rate with the given currency")
	meters = modProviderCmd.Flags().StringArray("meter", []string{}, "adds a meter with the given type to the preceeding rate")
	userTypes = modProviderCmd.Flags().StringArray("user-type", []string{}, "adds a user type to the preceeding rate")
	return modProviderCmd
}

func runModProviderCmd(cmd *cobra.Command, args []string) (err error) {
	if len(*rates) != len(*payments) || len(*meters) != len(*payments) || len(*userTypes) != len(*payments) {
		err = fmt.Errorf("must have equal number of payment, rate, meter, and user-types")
		return
	}

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

	subscriptionRates, payAsYouGoRates := cosmos.Coins{}, cosmos.Coins{}
	for i, payment := range *payments {
		rate, err := cosmos.ParseCoins(fmt.Sprintf("%su%s", (*rates)[i], payment))
		if err != nil {
			return err
		}
		meter := (*meters)[i]
		userType := (*userTypes)[i]
		_ = userType

		switch meter {
		case "block":
			subscriptionRates = subscriptionRates.Add(rate...)
		case "request":
			payAsYouGoRates = payAsYouGoRates.Add(rate...)
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
		subscriptionRates,
		payAsYouGoRates,
		int64(argSettlementDuration),
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
