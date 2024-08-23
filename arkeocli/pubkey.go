package arkeocli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

func newShowPubkeyCmd() *cobra.Command {
	bondProviderCmd := &cobra.Command{
		Use:   "show-pubkey",
		Short: "show pubkey for given key name or address (must exist in keyring)",
		Args:  cobra.ExactArgs(1),
		RunE:  runShowPubkeyCmd,
	}

	return bondProviderCmd
}

func runShowPubkeyCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return
	}
	addr, err := cosmos.AccAddressFromBech32(args[0])
	if err != nil {
		key, err := clientCtx.Keyring.Key(args[0])
		if err != nil {
			return err
		}
		addr, err = key.GetAddress()
		if err != nil {
			return err
		}
	}

	pubkey, err := toPubkey(cmd, addr)
	if err != nil {
		return
	}
	cmd.Println(pubkey)
	return
}
