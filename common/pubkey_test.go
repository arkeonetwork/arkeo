package common

import (
	"encoding/json"
	"testing"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
)

func TestPubKey(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	pk, err := NewPubKey(spk)
	require.NoError(t, err)
	hexStr := pk.String()
	require.True(t, len(hexStr) > 0)
	pk1, err := NewPubKey(hexStr)
	require.NoError(t, err)
	require.True(t, pk.Equals(pk1))

	result, err := json.Marshal(pk)
	require.NoError(t, err)

	var pk2 PubKey
	err = json.Unmarshal(result, &pk2)
	require.NoError(t, err)
	require.True(t, pk2.Equals(pk))
}

func TestEquals(t *testing.T) {
	var pk1, pk2, pk3, pk4 PubKey
	_, pubKey1, _ := testdata.KeyTestPubAddr()
	tpk1, err1 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
	require.NoError(t, err1)
	pk1 = PubKey(tpk1)

	_, pubKey2, _ := testdata.KeyTestPubAddr()
	tpk2, err2 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
	require.NoError(t, err2)
	pk2 = PubKey(tpk2)

	_, pubKey3, _ := testdata.KeyTestPubAddr()
	tpk3, err3 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
	require.NoError(t, err3)
	pk3 = PubKey(tpk3)

	_, pubKey4, _ := testdata.KeyTestPubAddr()
	tpk4, err4 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey4)
	require.NoError(t, err4)
	pk4 = PubKey(tpk4)

	require.False(t, PubKeys{pk1, pk2}.Equals(nil))
	require.False(t, PubKeys{pk1, pk2, pk3}.Equals(PubKeys{pk1, pk2}))
	require.True(t, PubKeys{pk1, pk2, pk3, pk4}.Equals(PubKeys{pk4, pk3, pk2, pk1}))
	require.True(t, PubKeys{pk1, pk2, pk3, pk4}.Equals(PubKeys{pk1, pk2, pk3, pk4})) // nolint
}
