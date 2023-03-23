package arkeocli

import (
	"bufio"
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
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

// find a contract by id or by provider/client/service
func getContract(cmd *cobra.Command) (*types.Contract, error) {
	var contract types.Contract
	queryCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return nil, err
	}

	queryClient := types.NewQueryClient(queryCtx)
	contractID, err := cmd.Flags().GetUint64("contract-id")
	if err != nil {
		return nil, err
	}
	if contractID != 0 && (cmd.Flags().Changed("provider-pubkey") || cmd.Flags().Changed("client-pubkey") || cmd.Flags().Changed("service")) {
		return nil, fmt.Errorf("cannot specify contract-id and provider-pubkey/client-pubkey/service")
	}

	if contractID != 0 {
		params := &types.QueryFetchContractRequest{ContractId: contractID}
		res, err := queryClient.FetchContract(cmd.Context(), params)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching contract %d", contractID)
		}
		if res == nil {
			return nil, fmt.Errorf("no contract found %d", contractID)
		}
		contract = res.GetContract()
	} else {
		providerPubkey, err := cmd.Flags().GetString("provider-pubkey")
		if err != nil {
			return nil, err
		}
		if providerPubkey == "" {
			providerPubkey, err = promptForArg(cmd, "Specify provider pubkey: ")
			if err != nil {
				return nil, err
			}
		}

		clientPubkey, err := cmd.Flags().GetString("client-pubkey")
		if err != nil {
			return nil, err
		}
		if clientPubkey == "" {
			clientPubkey, err = promptForArg(cmd, "Specify client pubkey: ")
			if err != nil {
				return nil, err
			}
		}

		service, err := cmd.Flags().GetString("service")
		if err != nil {
			return nil, err
		}
		if service == "" {
			service, err = promptForArg(cmd, "Specify service (e.g. gaia-mainnet-rpc-archive, btc-mainnet-fullnode, etc): ")
			if err != nil {
				return nil, err
			}
		}

		params := &types.QueryActiveContractRequest{
			Spender:  clientPubkey,
			Provider: providerPubkey,
			Service:  service,
		}

		res, err := queryClient.ActiveContract(cmd.Context(), params)
		if err != nil {
			return nil, errors.Wrapf(err, "could not find active contract for %s:%s:%s", clientPubkey, providerPubkey, service)
		}

		contract = res.GetContract()
	}
	return &contract, nil
}
