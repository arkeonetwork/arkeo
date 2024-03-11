package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/testutil/utils"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetClaimRecordForArkeo(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)

	addr1 := utils.GetRandomArkeoAddress().String()
	addr2 := utils.GetRandomArkeoAddress().String()
	addr3 := utils.GetRandomArkeoAddress().String()

	claimRecords := []types.ClaimRecord{
		{
			Chain:          types.ARKEO,
			Address:        addr1,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ARKEO,
			Address:        addr2,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		},
	}
	err := keepers.ClaimKeeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	coins1, err := keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr1, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, "300", coins1.Amount.String())

	coins2, err := keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr2, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, "600", coins2.Amount.String())

	coins3, err := keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr3, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins3, sdk.Coin{})

	claimRecord, err := keepers.ClaimKeeper.GetClaimRecord(ctx, addr3, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord, types.ClaimRecord{})

	// get rewards amount per action
	coins4, err := keepers.ClaimKeeper.GetClaimableAmountForAction(ctx, addr1, types.ACTION_VOTE, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins4.String(), sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)).String())
}

func TestGetClaimRecordForMutlipleChains(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)

	addr1 := utils.GetRandomArkeoAddress().String()
	addr2 := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address

	claimRecords := []types.ClaimRecord{
		{
			Chain:          types.ARKEO,
			Address:        addr1,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ETHEREUM,
			Address:        addr2,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		},
	}
	err := keepers.ClaimKeeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	coins1, err := keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr1, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, "300", coins1.Amount.String())

	// user 1 should have no eth claim with an arkeo addy nor thor claims
	coins1, err = keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr1, types.ETHEREUM)
	require.NoError(t, err)
	require.Equal(t, coins1, sdk.Coin{})

	// user 2 should have no arkeo claim nor thor claims, only eth
	coins2, err := keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr2, types.ETHEREUM)
	require.NoError(t, err)
	require.Equal(t, "600", coins2.Amount.String())

	coins2, err = keepers.ClaimKeeper.GetUserTotalClaimable(ctx, addr2, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins2, sdk.Coin{})
}

func TestSetClaimRecord(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)

	// confirm setting a claim record with a bad eth address fails
	addr1Invalid := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98"  // random invalid eth address
	addr1Valid := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random valid eth address
	claimRecord := types.ClaimRecord{
		Chain:          types.ETHEREUM,
		Address:        addr1Invalid,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err := keepers.ClaimKeeper.SetClaimRecord(ctx, claimRecord)
	require.Error(t, err)

	claimRecord.Address = addr1Valid
	err = keepers.ClaimKeeper.SetClaimRecord(ctx, claimRecord)
	require.NoError(t, err)

	// confirm setting a claim record with a bad arkeo address fails
	addr2Invalid := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address (should fail)
	addr2Valid := utils.GetRandomArkeoAddress().String()
	claimRecord = types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        addr2Invalid,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(ctx, claimRecord)
	require.Error(t, err)

	claimRecord.Address = addr2Valid
	err = keepers.ClaimKeeper.SetClaimRecord(ctx, claimRecord)
	require.NoError(t, err)
}

func TestGetAllClaimRecords(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)

	addr1 := utils.GetRandomArkeoAddress().String()
	addr2 := utils.GetRandomArkeoAddress().String()
	addr3 := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address

	claimRecords := []types.ClaimRecord{
		{
			Chain:          types.ARKEO,
			Address:        addr1,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ARKEO,
			Address:        addr2,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ETHEREUM,
			Address:        addr3,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		},
	}
	err := keepers.ClaimKeeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	// confirm all claims are returned
	claims, err := keepers.ClaimKeeper.GetAllClaimRecords(ctx)
	require.NoError(t, err)
	require.Equal(t, len(claims), len(claimRecords))
}

func TestClaimFlow(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	addrArkeo := utils.GetRandomArkeoAddress()

	claimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        addrArkeo.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err := keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimArkeo{
		Creator: addrArkeo.String(),
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, !claimRecord.IsEmpty())

	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// trigger event hook from voting
	keepers.ClaimKeeper.AfterProposalVote(sdkCtx, 1, addrArkeo)
	// confirm balance increased by expected amount.
	balanceAfter = keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// trigger another event from voting, nothing additional airdrop should happen.
	keepers.ClaimKeeper.AfterProposalVote(sdkCtx, 1, addrArkeo)
	balanceAfter = keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// trigger event hook from delegation
	err = keepers.ClaimKeeper.AfterDelegationModified(sdkCtx, addrArkeo, sdk.ValAddress(addrArkeo))
	require.NoError(t, err)
	// confirm balance increased by expected amount.
	balanceAfter = keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 300))

	// trigger event hook from voting, with an address that has no claim record
	addrArkeo2 := utils.GetRandomArkeoAddress()
	balanceBefore2 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo2, types.DefaultClaimDenom)
	keepers.ClaimKeeper.AfterProposalVote(sdkCtx, 1, addrArkeo2)
	balanceAfter2 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo2, types.DefaultClaimDenom)
	require.Equal(t, balanceBefore2, balanceAfter2)

	// same with delegation
	err = keepers.ClaimKeeper.AfterDelegationModified(sdkCtx, addrArkeo, sdk.ValAddress(addrArkeo))
	require.NoError(t, err)
	balanceAfter2 = keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo2, types.DefaultClaimDenom)
	require.Equal(t, balanceBefore2, balanceAfter2)
}

func TestClaimDecay(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	addrArkeo1 := utils.GetRandomArkeoAddress()
	addrArkeo2 := utils.GetRandomArkeoAddress()
	addrArkeo3 := utils.GetRandomArkeoAddress()

	claimRecords := []types.ClaimRecord{
		{
			Chain:          types.ARKEO,
			Address:        addrArkeo1.String(),
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ARKEO,
			Address:        addrArkeo2.String(),
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ARKEO,
			Address:        addrArkeo3.String(),
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
	}

	err := keepers.ClaimKeeper.SetClaimRecords(sdkCtx, claimRecords)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore1 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo1, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimArkeo{
		Creator: addrArkeo1.String(),
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.NoError(t, err)
	// trigger event hook from voting
	keepers.ClaimKeeper.AfterProposalVote(sdkCtx, 1, addrArkeo1)
	// trigger event hook from delegation
	err = keepers.ClaimKeeper.AfterDelegationModified(sdkCtx, addrArkeo1, sdk.ValAddress(addrArkeo1))
	require.NoError(t, err)

	// confirm balance increased by expected amount.
	balanceAfter1 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo1, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter1.Sub(balanceBefore1), sdk.NewInt64Coin(types.DefaultClaimDenom, 300))

	// advance time to middle of decay period
	params := keepers.ClaimKeeper.GetParams(sdkCtx)
	sdkCtx = sdkCtx.WithBlockTime(params.AirdropStartTime.Add(params.DurationUntilDecay).Add(params.DurationOfDecay / 2))

	// trigger all events for addr2
	balanceBefore2 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo2, types.DefaultClaimDenom)

	claimMessage = types.MsgClaimArkeo{
		Creator: addrArkeo2.String(),
	}
	_, err = msgServer.ClaimArkeo(sdkCtx, &claimMessage)
	require.NoError(t, err)
	// trigger event hook from voting
	keepers.ClaimKeeper.AfterProposalVote(sdkCtx, 1, addrArkeo2)
	// trigger event hook from delegation
	err = keepers.ClaimKeeper.AfterDelegationModified(sdkCtx, addrArkeo2, sdk.ValAddress(addrArkeo2))
	require.NoError(t, err)

	// confirm balance increased by expected amount. Should get half of original amounts!
	balanceAfter2 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo2, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter2.Sub(balanceBefore2), sdk.NewInt64Coin(types.DefaultClaimDenom, 150))

	// advance time to after end of airdrop
	sdkCtx = sdkCtx.WithBlockTime(params.AirdropStartTime.Add(params.DurationUntilDecay).Add(params.DurationOfDecay * 2))

	// trigger all events for addr3
	balanceBefore3 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo3, types.DefaultClaimDenom)

	claimMessage = types.MsgClaimArkeo{
		Creator: addrArkeo3.String(),
	}
	_, err = msgServer.ClaimArkeo(sdkCtx, &claimMessage)
	require.ErrorIs(t, err, types.ErrAirdropEnded)
	// trigger event hook from voting
	keepers.ClaimKeeper.AfterProposalVote(sdkCtx, 1, addrArkeo3)
	// trigger event hook from delegation
	err = keepers.ClaimKeeper.AfterDelegationModified(sdkCtx, addrArkeo3, sdk.ValAddress(addrArkeo3))
	require.NoError(t, err)

	// confirm balance increased by expected amount. Should get 0!
	balanceAfter3 := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo3, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter3.Sub(balanceBefore3), sdk.NewInt64Coin(types.DefaultClaimDenom, 0))
}
