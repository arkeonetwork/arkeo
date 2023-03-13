package arkeocli

import (
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var (
	bondProviderCmd = &cobra.Command{
		Use:   "bond-provider",
		Short: "bond or modify provider bond",
		Args:  cobra.ExactArgs(0),
		RunE:  runBondProviderCmd,
	}
)

func runBondProviderCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	fromAddr := clientCtx.GetFromAddress().String()
	if fromAddr == "" {
		readFrom, err := promptForArg(cmd, "Specify from key or address: ")
		if err != nil {
			return err
		}

		// readFrom := "alice"
		// readFrom := "tarkeo1up3pwhguqvr7l7pr8t53nmjrlxx03x0y0axw9z"

		var accAddr cosmos.AccAddress
		key, err := clientCtx.Keyring.Key(readFrom)
		if err != nil {
			accAddr, err = cosmos.AccAddressFromBech32(readFrom)
			if err != nil {
				return err
			}
			key, err = clientCtx.Keyring.KeyByAddress(accAddr)
			if err != nil {
				return err
			}
		}
		accAddr, err = key.GetAddress()
		if err != nil {
			return err
		}

		clientCtx = clientCtx.WithFromName(key.Name).WithFromAddress(accAddr)
		if err = client.SetCmdClientContext(cmd, clientCtx); err != nil {
			return err
		}
	}

	argPubkey, _ := cmd.Flags().GetString("pubkey")
	if argPubkey == "" {
		argPubkey, err = promptForArg(cmd, "Specify provider pubkey: ")
		if err != nil {
			return err
		}
	}

	argChain, _ := cmd.Flags().GetString("chain")
	if argChain == "" {
		argChain, err = promptForArg(cmd, "Specify chain (e.g. gaia-mainnet-rpc-archive, btc-mainnet-fullnode, etc): ")
		if err != nil {
			return err
		}
	}

	// ensure valid chain
	_, err = common.NewChain(argChain)
	if err != nil {
		return
	}

	argBond, _ := cmd.Flags().GetString("bond")
	if argBond == "" {
		argBond, err = promptForArg(cmd, "Specify bond amount (e.g. 100uarkeo): ")
		if err != nil {
			return err
		}
	}

	coins, err := cosmos.ParseCoins(argBond)
	if err != nil {
		return err
	}
	if len(coins) != 1 {
		return fmt.Errorf("1 coins as bond amount, got %d", len(coins))
	}
	if coins[0].Denom != "uarkeo" {
		return fmt.Errorf("bad bond denom, expected \"uarkeo\" got \"%s\"", coins[0].Denom)
	}
	if coins[0].Amount.IsNegative() || coins[0].Amount.IsZero() {
		return fmt.Errorf("bad bond amount: %s", argBond)
	}

	pubkey, err := common.NewPubKey(argPubkey)
	if err != nil {
		return err
	}
	msg := types.NewMsgBondProvider(
		clientCtx.GetFromAddress().String(),
		pubkey,
		argChain,
		coins[0].Amount,
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
