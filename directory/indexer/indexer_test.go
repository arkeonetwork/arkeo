package indexer

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

const (
	bech32PrefixAccAddr = "arkeo"
	bech32PrefixAccPub  = "arkeopub"
)

// list keys in $HOME/.arkeo/keyring-test with addr and pubkey
func TestAccountDetails(t *testing.T) {
	arkeoHome := fmt.Sprintf("%s/.arkeo", os.Getenv("HOME"))
	err := accountDetails(arkeoHome)
	if err != nil {
		t.Errorf("error getting details for alice: %+v", err)
	}
}

func accountDetails(keyringPath string) error {
	var (
		addr, pubkey string
		err          error
	)
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	encConfig := NewEncoding()

	keyRing, err := keyring.New("arkeo", "test", keyringPath, nil, encConfig.Marshaler)
	if err != nil {
		return err
	}
	all, err := keyRing.List()
	if err != nil {
		return err
	}
	for _, v := range all {
		pub, perr := v.GetPubKey()
		if perr != nil {
			log.Errorf("error getting \"%s\" pubkey from keyring: %+v", v.Name, err)
			continue
		}
		accAddr := sdk.AccAddress(pub.Address())
		addr, err = bech32.ConvertAndEncode(sdkConfig.GetBech32AccountAddrPrefix(), accAddr)
		if err != nil {
			log.Errorf("error encoding account address %+v", err)
			continue
		}
		pubkey, err = bech32.ConvertAndEncode(sdkConfig.GetBech32AccountPubPrefix(), legacy.Cdc.MustMarshal(pub))
		if err != nil {
			log.Errorf("error encoding pubkey %+v", err)
			return err
		}

		log.Infof("%s addr: %s pubkey: %s", v.Name, addr, pubkey)
	}
	return nil
}
