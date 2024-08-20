package types

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestClaimContractIncomeValidateBasic(t *testing.T) {
	// setup
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)
	kb := cKeys.NewInMemory(cdc)
	_, _, err = kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	require.NoError(t, err)

	msg := MsgClaimContractIncome{
		Creator:    acct.String(),
		ContractId: 1,
		Nonce:      24,
	}

	message := msg.GetBytesToSign()
	msg.Signature, _, err = kb.Sign("whatever", message, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(t, err)
	err = msg.ValidateBasic()
	require.NoError(t, err)
}

func TestValidateSignature(t *testing.T) {
	t.Skip("kept here for archival / 'working with pubkeys' purposes")
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	msg := MsgClaimContractIncome{
		Creator:    acct.String(),
		Nonce:      48,
		ContractId: 500,
	}
	err = msg.ValidateBasic()
	require.NoError(t, err)

	// setup
	var pub cryptotypes.PubKey
	kb := cKeys.NewInMemory(cdc)
	_, _, err = kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	require.NoError(t, err)

	message := []byte(fmt.Sprintf("%d:%d", msg.ContractId, msg.Nonce))
	msg.Signature, pub, err = kb.Sign("whatever", message, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(t, err)

	require.True(t, pub.VerifySignature(message, msg.Signature))

	pk, err := common.NewPubKeyFromCrypto(pub)
	require.NoError(t, err)

	pk2, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, pk.String())
	require.NoError(t, err)
	require.True(t, pk2.Equals(pub))

	acc, err := pk.GetMyAddress()
	require.NoError(t, err)
	account := authtypes.NewBaseAccountWithAddress(acc)
	require.True(t, pk2.Equals(account.GetPubKey()))
}
