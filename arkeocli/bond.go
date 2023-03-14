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
	bondProviderCmd.Flags().StringP("pubkey", "p", "", "provider pubkey")
	bondProviderCmd.Flags().StringP("chain", "c", "", "provider chain")
	bondProviderCmd.Flags().String("bond", "", "provider bond amount")
	return bondProviderCmd
}

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
		argBond, err = promptForArg(cmd, "Specify bond amount (e.g. 100uarkeo, negative to unbond): ")
		if err != nil {
			return err
		}
	}

	// cosmos.ParseCoins fails negative numbers
	bondIntArr := make([]rune, 0, len(argBond))
	var i int
	var c rune
	for i, c = range argBond {
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
	bondDenom := argBond[i:]
	if bondDenom != "uarkeo" {
		return fmt.Errorf("bad bond denom, expected \"uarkeo\" got \"%s\"", bondDenom)
	}

	bond, ok := cosmos.NewIntFromString(string(bondIntArr))
	if !ok {
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
		bond,
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
