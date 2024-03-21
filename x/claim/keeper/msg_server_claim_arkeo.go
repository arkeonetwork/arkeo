package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

type ThorTxData struct {
	ObservedTx struct {
		Tx struct {
			FromAddress string `json:"from_address"`
		} `json:"tx"`
	} `json:"observed_tx"`
}

func (k msgServer) ClaimArkeo(goCtx context.Context, msg *types.MsgClaimArkeo) (*types.MsgClaimArkeoResponse, error) {
	log.Println("WHAT")
	ctx := sdk.UnwrapSDKContext(goCtx)
	arkeoClaim, err := k.GetClaimRecord(ctx, msg.Creator.String(), types.ARKEO)
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

		thorClaim, err := k.GetClaimRecord(ctx, thorDerivedArkeoAddress, types.ARKEO)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get claim record for %s", thorDerivedArkeoAddress)
		}
		if thorClaim.IsEmpty() || thorClaim.AmountClaim.IsZero() {
			return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", thorDerivedArkeoAddress)
		}
		log.Println("Thor Claim: ", thorClaim)

		// TODO: Update claim record for arkeo address and remove claim for thor address
	}

	if arkeoClaim.IsEmpty() || arkeoClaim.AmountClaim.IsZero() {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.Creator)
	}

	_, err = k.ClaimCoinsForAction(ctx, msg.Creator.String(), types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.Creator)
	}

	return &types.MsgClaimArkeoResponse{}, nil
}
