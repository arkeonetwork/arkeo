package arkeocli

import (
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func newBondProviderCmd() *cobra.Command {
	bondProviderCmd := &cobra.Command{
		Use:   "bond-provider",
		Short: "bond or modify provider bond",
		Args:  cobra.ExactArgs(0),
		RunE:  runBondProviderCmd,
	}

	flags.AddTxFlagsToCmd(bondProviderCmd)
	bondProviderCmd.Flags().String("provider-pubkey", "", "provider pubkey")
	bondProviderCmd.Flags().String("service", "", "provider service name")
	bondProviderCmd.Flags().String("bond", "", "provider bond amount (negative to unbond)")
	return bondProviderCmd
}

func runBondProviderCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	key, err := ensureKeys(cmd)
	if err != nil {
		return err
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
		argPubkey, err = toPubKey(cmd, addr)
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

	argBond, _ := cmd.Flags().GetString("bond")
	if argBond == "" {
		argBond, err = promptForArg(cmd, "Specify bond amount (e.g. 100uarkeo, negative to unbond): ")
		if err != nil {
			return err
		}
	}
	bond, err := parseBondAmount(argBond)
	if err != nil {
		return err
	}

	pubkey, err := common.NewPubKey(argPubkey)
	if err != nil {
		return err
	}
	msg := types.NewMsgBondProvider(
		clientCtx.GetFromAddress(),
		pubkey,
		argService,
		bond,
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}

func parseBondAmount(bondStr string) (amount cosmos.Int, err error) {
	// cosmos.ParseCoins fails negative numbers
	bondIntArr := make([]rune, 0, len(bondStr))
	var i int
	var c rune
	for i, c = range bondStr {
		if i == 0 && c == '-' {
			bondIntArr = append(bondIntArr, c)
			continue
		}
		if c >= '0' && c <= '9' {
			bondIntArr = append(bondIntArr, c)
			continue
		}
		break
	}
	bondDenom := bondStr[i:]
	if bondDenom != "uarkeo" {
		err = fmt.Errorf("bad bond denom, expected \"uarkeo\" got \"%s\"", bondDenom)
		return
	}

	var ok bool
	amount, ok = cosmos.NewIntFromString(string(bondIntArr))
	if !ok {
		err = fmt.Errorf("bad bond amount: %s", bondStr)
		return
	}

	return
}
