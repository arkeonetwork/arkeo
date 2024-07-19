package keeper

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	crypto "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/pkg/errors"
)

// Verify and update the claim record based on the memo of the thorchain tx
func (k msgServer) updateThorClaimRecord(ctx sdk.Context, creator string, thorTxMsg *types.MsgThorTxData, arkeoClaimRecord types.ClaimRecord) (types.ClaimRecord, error) {

	thorTxChainData, thorTxEncodedData, err := decodeAndUnmarshallThorTxMsg(thorTxMsg)
	if err != nil {
		return types.ClaimRecord{}, err
	}

	verifyTxDataHash, txDataHex, err := verifyTxDataCheckSum(thorTxChainData, thorTxEncodedData)
	if err != nil {
		return types.ClaimRecord{}, err
	}

	err = verifySignature(thorTxMsg, verifyTxDataHash, txDataHex)

	if err != nil {
		return types.ClaimRecord{}, err
	}

	thorAddress := thorTxChainData.ObservedTx.Tx.FromAddress
	memo := thorTxChainData.ObservedTx.Tx.Memo

	thorAddressBytes, err := sdk.GetFromBech32(thorAddress, "thor")
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("invalid thorchain address: %w", err)
	}
	prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

	// Re-encode the raw bytes with the new prefix
	thorDerivedArkeoAddress, err := sdk.Bech32ifyAddressBytes(prefix, thorAddressBytes)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("failed to encode address bytes with new prefix: %w", err)
	}

	thorClaimRecord, err := k.GetClaimRecord(ctx, thorDerivedArkeoAddress, types.ARKEO)
	if err != nil {
		return types.ClaimRecord{}, errors.Wrapf(err, "failed to get claim record for %s", thorDerivedArkeoAddress)
	}
	if thorClaimRecord.IsEmpty() || thorClaimRecord.AmountClaim.IsZero() {
		return types.ClaimRecord{}, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", thorDerivedArkeoAddress)
	}
	parts := strings.Split(memo, ":")
	if len(parts) != 3 || parts[0] != "delegate" || parts[1] != "arkeo" {
		return types.ClaimRecord{}, fmt.Errorf("invalid memo: %s", memo)
	}

	combinedClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        creator,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, arkeoClaimRecord.AmountClaim.Amount.Int64()+thorClaimRecord.AmountClaim.Amount.Int64()),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, arkeoClaimRecord.AmountVote.Amount.Int64()+thorClaimRecord.AmountVote.Amount.Int64()),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, arkeoClaimRecord.AmountDelegate.Amount.Int64()+thorClaimRecord.AmountDelegate.Amount.Int64()),
	}
	emptyClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        thorDerivedArkeoAddress,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
	}
	err = k.SetClaimRecord(ctx, emptyClaimRecord)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("failed to set empty claim record for %s: %w", thorDerivedArkeoAddress, err)
	}
	err = k.SetClaimRecord(ctx, combinedClaimRecord)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("failed to set combined claim record for %s: %w", creator, err)
	}

	newClaimRecord, err := k.GetClaimRecord(ctx, creator, types.ARKEO)
	if err != nil {
		return types.ClaimRecord{}, errors.Wrapf(err, "failed to get claim record for %s", creator)
	}
	return newClaimRecord, nil
}

func decodeAndUnmarshallThorTxMsg(thorTxMsg *types.MsgThorTxData) (*types.ThorChainTxData, *types.ThorTxData, error) {
	// decode thorTxMessage
	thorTxMsgDecoded, err := hex.DecodeString(thorTxMsg.ThorData)
	if err != nil {
		return nil, nil, fmt.Errorf("error hex decoding faild: %w", err)
	}

	// unmarshall encoded data
	var thorTxEncodedData *types.ThorTxData
	if err := json.Unmarshal(thorTxMsgDecoded, &thorTxEncodedData); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling transaction data: %w", err)
	}

	// decode tx data
	thorTxDataDecoded, err := hex.DecodeString(thorTxEncodedData.TxData)
	if err != nil {
		return nil, nil, fmt.Errorf("error hex decoding failed: %w", err)
	}

	// unmarshall thor tx
	var thorTxChainData *types.ThorChainTxData
	if err := json.Unmarshal(thorTxDataDecoded, &thorTxChainData); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling transaction data: %w", err)
	}

	return thorTxChainData, thorTxEncodedData, nil
}

func verifyTxDataCheckSum(thorTxChainData *types.ThorChainTxData, thorTxEncodedData *types.ThorTxData) (string, string, error) {
	thorTxChainDataBytes, err := json.Marshal(thorTxChainData)
	if err != nil {
		return "", "", fmt.Errorf("error marshalling tx data: %w", err)
	}

	txDataHash := sha512.Sum512(thorTxChainDataBytes)
	verifyTxDataHash := hex.EncodeToString(txDataHash[:])
	txDataHex := hex.EncodeToString(thorTxChainDataBytes)

	// Verify Data Check Sum
	if verifyTxDataHash != thorTxEncodedData.Hash {
		return "", "", fmt.Errorf("transaction data cehcksum failed")
	}

	return verifyTxDataHash, txDataHex, nil

}

func verifySignature(thorTxMsg *types.MsgThorTxData, verifyTxDataHash, txDataHex string) error {
	txData := types.ThorTxData{
		Hash:   verifyTxDataHash,
		TxData: txDataHex,
	}

	txDataBytes, err := json.Marshal(txData)
	if err != nil {
		return fmt.Errorf("error marshalling txData: %w", err)
	}

	txDataHashForVerification := hex.EncodeToString(txDataBytes)

	// verify the tx signature
	proofPubKeyDecoded, err := hex.DecodeString(thorTxMsg.ProofPubkey)
	if err != nil {
		return fmt.Errorf("failed to decode pub key: %w", err)
	}
	pubkey := crypto.PubKey{}
	pubkey.Key = proofPubKeyDecoded

	proofSignatureDecoded, err := hex.DecodeString(thorTxMsg.ProofSignature)
	if err != nil {
		return fmt.Errorf("failed to decode signature key: %w", err)
	}

	// create data to verify
	hash := sha512.Sum512([]byte(txDataHashForVerification))

	if !pubkey.VerifySignature(hash[:], proofSignatureDecoded) {
		return fmt.Errorf("message verification failed")
	}

	return nil
}
