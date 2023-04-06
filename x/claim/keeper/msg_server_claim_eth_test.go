package keeper_test

import (
	"crypto/ecdsa"
	"errors"
	"strings"
	"testing"

	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestClaimEth(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo.String(), "300")
	require.NoError(t, err)

	claimRecord := types.ClaimRecord{
		Chain:          types.ETHEREUM,
		Address:        addrEth,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimEth{
		Creator:    addrArkeo,
		EthAddress: addrEth,
		Signature:  sigString,
	}
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrEth, types.ETHEREUM)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.Error(t, err)

	// attempt to claim from arkeo should also fail!
	_, err = msgServer.ClaimArkeo(ctx, &types.MsgClaimArkeo{Creator: addrArkeo})
	require.Error(t, err)
}

func TestClaimEthWithInvalidSignature(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo.String(), "200")
	require.NoError(t, err)

	claimRecord := types.ClaimRecord{
		Chain:          types.ETHEREUM,
		Address:        addrEth,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	claimMessage := types.MsgClaimEth{
		Creator:    addrArkeo,
		EthAddress: addrEth,
		Signature:  sigString,
	}
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.ErrorIs(t, types.ErrInvalidSignature, err)
}

func TestClaimEthWithArkeoClaimRecord(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo.String(), "600")
	require.NoError(t, err)

	claimRecord := types.ClaimRecord{
		Chain:          types.ETHEREUM,
		Address:        addrEth,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	// create an arkeo claim record for the same user. This should be merged once they call claim.
	claimRecordArkeo := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        addrArkeo.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		AmountVote:     sdk.Coin{},
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 150),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecordArkeo)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimEth{
		Creator:    addrArkeo,
		EthAddress: addrEth,
		Signature:  sigString,
	}
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrEth, types.ETHEREUM)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 350))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 400))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.Error(t, err)

	// attempt to claim from arkeo should also fail!
	_, err = msgServer.ClaimArkeo(ctx, &types.MsgClaimArkeo{Creator: addrArkeo})
	require.Error(t, err)
}

func TestClaimEthWithNoClaimRecord(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo.String(), "300")
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	claimMessage := types.MsgClaimEth{
		Creator:    addrArkeo,
		EthAddress: addrEth,
		Signature:  sigString,
	}
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.Error(t, err)
}

func TestIsValidClaimSignature(t *testing.T) {
	// generate a random eth address
	addrArkeo := utils.GetRandomArkeoAddress().String()
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
	addrArkeo2 := utils.GetRandomArkeoAddress().String()
	_, err = keeper.IsValidClaimSignature(addressEth, addrArkeo2, "5000", sigString)
	require.Error(t, err)

	// if we modify the eth address, signature should be invalid
	_, err = keeper.IsValidClaimSignature("0xbd3afb0bb76683ecb4225f9dbc91f998713c3b01", addrArkeo, "5000", sigString)
	require.Error(t, err)
}

func generateSignedEthClaim(addrArkeo, amount string) (string, string, error) {
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
