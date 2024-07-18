package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// Verify and update the claim record based on the memo of the thorchain tx
func (k msgServer) updateThorClaimRecord(ctx sdk.Context, creator string, thorTxMsg *types.MsgThorTxData, arkeoClaimRecord types.ClaimRecord) (types.ClaimRecord, error) {

	thorDataDecoded, err := base64.StdEncoding.DecodeString(thorTxMsg.ThorData)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("error base64 decoding faild: %w", err)
	}
	var thorData *types.ThorTxData
	if err := json.Unmarshal(thorDataDecoded, &thorData); err != nil {
		return types.ClaimRecord{}, fmt.Errorf("error unmarshalling transaction data: %w", err)
	}

	var thorTxData *types.ThorChainTxData
	if err := json.Unmarshal([]byte(thorData.TxData), &thorTxData); err != nil {
		return types.ClaimRecord{}, fmt.Errorf("error unmarshalling transaction data: %w", err)
	}

	marshalledtxBytes, err := json.Marshal(thorData)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("error marshalling tx data: %w", err)
	}

	verifyTxData := base64.StdEncoding.EncodeToString(marshalledtxBytes)

	if verifyTxData != thorData.Hash {
		return types.ClaimRecord{}, fmt.Errorf("transaction data cehcksum failed")
	}

	thorAddress := thorTxData.ObservedTx.Tx.FromAddress
	memo := thorTxData.ObservedTx.Tx.Memo

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

	decodedPubKey, err := hex.DecodeString(thorTxMsg.ProofPubkey)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("error in decoding pubkey: %w", err)
	}

	proofAddress, err := sdk.Bech32ifyAddressBytes(prefix, decodedPubKey)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("invalid thorchain address: %w", err)
	}

	if proofAddress != thorDerivedArkeoAddress {
		return types.ClaimRecord{}, fmt.Errorf("address validation failed: %w", err)
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
