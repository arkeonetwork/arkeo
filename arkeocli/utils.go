package arkeocli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/spf13/cobra"
)

func promptForArg(cmd *cobra.Command, prompt string) (string, error) {
	cmd.Print(prompt)
	reader := bufio.NewReader(cmd.InOrStdin())
	read, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	read = strings.TrimSpace(read)
	return read, nil
}

func toPubkey(cmd *cobra.Command, addr cosmos.AccAddress) (pubkey string, err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return
	}
	key, err := clientCtx.Keyring.KeyByAddress(addr)
	if err != nil {
		return
	}
	pub, err := key.GetPubKey()
	if err != nil {
		return
	}
	pubkey, err = bech32.ConvertAndEncode(getAcctPubPrefix(), legacy.Cdc.MustMarshal(pub))
	if err != nil {
		return
	}

	return pubkey, nil
}

func getAcctPubPrefix() string {
	prefix := fmt.Sprintf("%spub", app.AccountAddressPrefix)
	return prefix
}

func ensureKeys(cmd *cobra.Command) (key *keyring.Record, err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return
	}
	fromKeyName := clientCtx.GetFromName()
	fromAddr := clientCtx.GetFromAddress()
	if fromKeyName != "" && !fromAddr.Empty() {
		key, err = clientCtx.Keyring.Key(fromKeyName)
		if err != nil {
			return
		}
		return key, nil
	}

	readFrom, err := promptForArg(cmd, "Specify from key or address: ")
	if err != nil {
		return
	}

	// try as key id
	key, err = clientCtx.Keyring.Key(readFrom)
	if err != nil {
		// try as bech32 address
		fromAddr, err = cosmos.AccAddressFromBech32(readFrom)
		if err != nil {
			return
		}
		key, err = clientCtx.Keyring.KeyByAddress(fromAddr)
		if err != nil {
			return
		}
	}

	return
}
