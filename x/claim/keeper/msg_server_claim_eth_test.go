package keeper_test

import (
	"crypto/ecdsa"
	"errors"
	"strings"
	"testing"

	"github.com/arkeonetwork/arkeo/x/claim/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestClaimEth(t *testing.T) {
	msgServer, keeper, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo, "100")
	require.NoError(t, err)

	claimRecord := types.ClaimRecord{
		Chain:                  types.ETHEREUM,
		Address:                addrEth,
		InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
		ActionCompleted:        []bool{false, false},
	}
	err = keeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	claimMessage := types.MsgClaimEth{
		Creator:    addrArkeo,
		EthAddress: addrEth,
		Signature:  sigString,
	}

	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keeper.GetClaimRecord(sdkCtx, addrEth, types.ETHEREUM)
	require.NoError(t, err)
	require.True(t, claimRecord.ActionCompleted[types.FOREIGN_CHAIN_ACTION_CLAIM])

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keeper.GetClaimRecord(sdkCtx, addrArkeo, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, addrArkeo)
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.Equal(t, claimRecord.InitialClaimableAmount, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)))
	require.False(t, claimRecord.ActionCompleted[types.ACTION_VOTE])
	require.False(t, claimRecord.ActionCompleted[types.ACTION_DELEGATE_STAKE])
}

func TestIsValidClaimSignature(t *testing.T) {
	// generate a random eth address
	addrArkeo := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addressEth, sigString, err := generateSignedEthClaim(addrArkeo, "5000")
	require.NoError(t, err)

	// check if signature is valid
	valid, err := keeper.IsValidClaimSignature(strings.ToLower(addressEth), addrArkeo, "5000", sigString)
	require.NoError(t, err)
	require.True(t, valid)

	// if we modify the message, signature should be invalid
	_, err = keeper.IsValidClaimSignature(addressEth, addrArkeo, "5001", sigString)
	require.Error(t, err)

	// if we modify the arkeo address, signature should be invalid
	addrArkeo2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	_, err = keeper.IsValidClaimSignature(addressEth, addrArkeo2, "5000", sigString)
	require.Error(t, err)

	// if we modify the eth address, signature should be invalid
	_, err = keeper.IsValidClaimSignature("0xbd3afb0bb76683ecb4225f9dbc91f998713c3b01", addrArkeo, "5000", sigString)
	require.Error(t, err)
}

func generateSignedEthClaim(addrArkeo string, amount string) (string, string, error) {
	// generate a random eth address
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("error casting public key to ECDSA")
	}

	addressEth := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	message, err := keeper.GenerateClaimTypedDataBytes(addressEth, addrArkeo, amount)
	if err != nil {
		return "", "", err
	}
	hash := crypto.Keccak256(message)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return "", "", err
	}
	sigString := hexutil.Encode(signature)
	return addressEth, sigString, nil
}
