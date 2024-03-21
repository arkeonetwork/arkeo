package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type ThorTxData struct {
	ObservedTx struct {
		Tx struct {
			FromAddress string `json:"from_address"`
			Memo        string `json:"memo"`
		} `json:"tx"`
	} `json:"observed_tx"`
}

func (k msgServer) ClaimArkeo(goCtx context.Context, msg *types.MsgClaimArkeo) (*types.MsgClaimArkeoResponse, error) {
	log.Println("WHAT")
	ctx := sdk.UnwrapSDKContext(goCtx)
	arkeoClaimRecord, err := k.GetClaimRecord(ctx, msg.Creator.String(), types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.Creator)
	}

	if msg.ThorTx != "" {
		log.Println("Thor Tx: ", msg.ThorTx)
		url := fmt.Sprintf("https://thornode.ninerealms.com/thorchain/tx/%s", msg.ThorTx)

		resp, err := http.Get(url)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get thorchain tx for %s", msg.ThorTx)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		var txData ThorTxData
		if err := json.Unmarshal(body, &txData); err != nil {
			return nil, fmt.Errorf("error unmarshalling transaction data: %w", err)
		}
		thorAddress := txData.ObservedTx.Tx.FromAddress
		memo := txData.ObservedTx.Tx.Memo

		thorAddressBytes, err := sdk.GetFromBech32(thorAddress, "thor")
		if err != nil {
			// thorAddress is invalid
			return nil, fmt.Errorf("not a thor tx: %w", err)
		}
		prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()

		// Re-encode the raw bytes with the new prefix
		thorDerivedArkeoAddress, err := sdk.Bech32ifyAddressBytes(prefix, thorAddressBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to encode address bytes with new prefix: %w", err)
		}

		thorClaimRecord, err := k.GetClaimRecord(ctx, thorDerivedArkeoAddress, types.ARKEO)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get claim record for %s", thorDerivedArkeoAddress)
		}
		if thorClaimRecord.IsEmpty() || thorClaimRecord.AmountClaim.IsZero() {
			return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", thorDerivedArkeoAddress)
		}
		parts := strings.Split(memo, ":")
		if len(parts) != 3 || parts[0] != "delegate" || parts[1] != "arkeo" {
			return nil, fmt.Errorf("invalid memo: %s", memo)
		}

		combinedClaimRecord := types.ClaimRecord{
			Chain:          types.ARKEO,
			Address:        msg.Creator.String(),
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
		k.SetClaimRecord(ctx, emptyClaimRecord)
		k.SetClaimRecord(ctx, combinedClaimRecord)

		thorClaimRecord, _ = k.GetClaimRecord(ctx, thorDerivedArkeoAddress, types.ARKEO)
		arkeoClaimRecord, err = k.GetClaimRecord(ctx, msg.Creator.String(), types.ARKEO)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.Creator)
		}
	}
	log.Println("Arkeo Claim: ", arkeoClaimRecord)
	if arkeoClaimRecord.IsEmpty() || arkeoClaimRecord.AmountClaim.IsZero() {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.Creator)
	}

	_, err = k.ClaimCoinsForAction(ctx, msg.Creator.String(), types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.Creator)
	}

	return &types.MsgClaimArkeoResponse{}, nil
}
