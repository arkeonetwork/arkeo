package keeper_test

import (
	"testing"

	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestClaimThorchainArkeo(t *testing.T) {
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

	thorClaimAddress := "cosmos1dllfyp57l4xj5umqfcqy6c2l3xfk0qk6wy5w8c"
	thorClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        thorClaimAddress, // arkeo address derived from sender of thorchain tx "FA2768AEB52AE0A378372B48B10C5B374B25E8B2005C702AAD441B813ED2F174"
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, thorClaimRecord)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimArkeo{
		Creator: addrArkeo,
		ThorTxData: &types.MsgThorTxData{
			ThorData:       "7b2268617368223a223133373430646435623638613938356662386364333464323737353230373039643637653065623939633665356631663430333036366233393662336566656533306138663765653539303165663432313535393036636561626439356538393834323132353439643235336536303034346133366361643934346538383835222c2274785f64617461223a223762323236663632373336353732373636353634356637343738323233613762323237343738323233613762323236393634323233613232343634313332333733363338343134353432333533323431343533303431333333373338333333373332343233343338343233313330343333353432333333373334343233323335343533383432333233303330333534333337333033323431343134343334333433313432333833313333343534343332343633313337333432323263323236333638363136393665323233613232353434383466353232323263323236363732366636643566363136343634373236353733373332323361323237343638366637323331363436633663363637393730333533373663333437383661333537353664373136363633373137393336363333323663333337383636366233303731366233363637373236343334366133383232326332323734366635663631363436343732363537333733323233613232373436383666373233313637333933383633373933333665333936643664366137323730366533303733373836643665333633333663376137343635366336353732363133333337366533383665333633373633333032323263323236333666363936653733323233613562376232323631373337333635373432323361323235343438346635323265353235353465343532323263323236313664366637353665373432323361323233303232376435643263323236373631373332323361366537353663366332633232366436353664366632323361323236343635366336353637363137343635336136313732366236353666336137343631373236623635366633313339333333353338376133323336366137373638333336353334373236343336373037333738373136363338373133363636333337303635333636363338373333373736333037383332363132323764376432633232363336663665373336353665373337353733356636383635363936373638373432323361333133353331333833333334333333303263323236363639366536313663363937333635363435663638363536393637363837343232336133313335333133383333333433333330326332323662363537393733363936373665356636643635373437323639363332323361376232323734373835663639363432323361323234363431333233373336333834313435343233353332343134353330343133333337333833333337333234323334333834323331333034333335343233333337333434323332333534353338343233323330333033353433333733303332343134313434333433343331343233383331333334353434333234363331333733343232326332323665366636343635356637343733373335663734363936643635373332323361366537353663366337643764227d",
			ProofSignature: "8af1915a046a5b3a11a1c4bf5f8f30f6e05a590a1b3361f69ee8797dd4e6a3ad7679d7fcf359c500cf71d645a215c888ab3e39b8082b2c5975ad5ed8d5004c44",
			ProofPubkey:    "020C8FF0D34D4B5EE779540ECA039B1DAC0F8EDFE9BE6EC233AA4B4FF8DE6EBF86",
		},
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	thorClaimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, thorClaimAddress, types.ARKEO)
	require.NoError(t, err)
	require.True(t, thorClaimRecord.IsEmpty())

	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, !claimRecord.IsEmpty())

	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.ErrorIs(t, err, types.ErrNoClaimableAmount)

	// ensure claim Arkeo fails from address with no claim record
	addrArkeo2 := utils.GetRandomArkeoAddress()
	claimMessage2 := types.MsgClaimArkeo{
		Creator: addrArkeo2,
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage2)
	require.ErrorIs(t, err, types.ErrNoClaimableAmount)
}

func TestClaimThorchainEth(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo.String(), "300")
	require.NoError(t, err)

	arkeoClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        addrArkeo.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 50),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 50),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 50),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, arkeoClaimRecord)
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

	thorClaimAddress := "cosmos1dllfyp57l4xj5umqfcqy6c2l3xfk0qk6wy5w8c"
	thorClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        thorClaimAddress, // arkeo address derived from sender of thorchain tx "FA2768AEB52AE0A378372B48B10C5B374B25E8B2005C702AAD441B813ED2F174"
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 500),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 500),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 500),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, thorClaimRecord)
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
		ThorTx: &types.MsgThorTxData{
			ThorData:       "7b2268617368223a223133373430646435623638613938356662386364333464323737353230373039643637653065623939633665356631663430333036366233393662336566656533306138663765653539303165663432313535393036636561626439356538393834323132353439643235336536303034346133366361643934346538383835222c2274785f64617461223a223762323236663632373336353732373636353634356637343738323233613762323237343738323233613762323236393634323233613232343634313332333733363338343134353432333533323431343533303431333333373338333333373332343233343338343233313330343333353432333333373334343233323335343533383432333233303330333534333337333033323431343134343334333433313432333833313333343534343332343633313337333432323263323236333638363136393665323233613232353434383466353232323263323236363732366636643566363136343634373236353733373332323361323237343638366637323331363436633663363637393730333533373663333437383661333537353664373136363633373137393336363333323663333337383636366233303731366233363637373236343334366133383232326332323734366635663631363436343732363537333733323233613232373436383666373233313637333933383633373933333665333936643664366137323730366533303733373836643665333633333663376137343635366336353732363133333337366533383665333633373633333032323263323236333666363936653733323233613562376232323631373337333635373432323361323235343438346635323265353235353465343532323263323236313664366637353665373432323361323233303232376435643263323236373631373332323361366537353663366332633232366436353664366632323361323236343635366336353637363137343635336136313732366236353666336137343631373236623635366633313339333333353338376133323336366137373638333336353334373236343336373037333738373136363338373133363636333337303635333636363338373333373736333037383332363132323764376432633232363336663665373336353665373337353733356636383635363936373638373432323361333133353331333833333334333333303263323236363639366536313663363937333635363435663638363536393637363837343232336133313335333133383333333433333330326332323662363537393733363936373665356636643635373437323639363332323361376232323734373835663639363432323361323234363431333233373336333834313435343233353332343134353330343133333337333833333337333234323334333834323331333034333335343233333337333434323332333534353338343233323330333033353433333733303332343134313434333433343331343233383331333334353434333234363331333733343232326332323665366636343635356637343733373335663734363936643635373332323361366537353663366337643764227d",
			ProofSignature: "8af1915a046a5b3a11a1c4bf5f8f30f6e05a590a1b3361f69ee8797dd4e6a3ad7679d7fcf359c500cf71d645a215c888ab3e39b8082b2c5975ad5ed8d5004c44",
			ProofPubkey:    "020C8FF0D34D4B5EE779540ECA039B1DAC0F8EDFE9BE6EC233AA4B4FF8DE6EBF86",
		},
	}

	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrEth, types.ETHEREUM)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	thorClaimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, thorClaimAddress, types.ARKEO)
	require.NoError(t, err)
	require.True(t, thorClaimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 650))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 650))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 650))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.Error(t, err)

	// attempt to claim from arkeo should also fail!
	_, err = msgServer.ClaimArkeo(ctx, &types.MsgClaimArkeo{Creator: addrArkeo})
	require.Error(t, err)
}
