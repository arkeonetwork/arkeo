package types

import (
	"math/rand"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/tendermint/tendermint/crypto"
)

const DefaultSignatureExpiration = int64(50) // SIGNATURE EXPIRES AFTER 50 Blocks

// GetRandomBech32Addr is an account address used for test
func GetRandomBech32Addr() cosmos.AccAddress {
	name := common.RandStringBytesMask(10)
	return cosmos.AccAddress(crypto.AddressHash([]byte(name)))
}

func GetRandomPubKey() common.PubKey {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	accts := simtypes.RandomAccounts(r, 1)
	bech32PubKey, _ := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, accts[0].PubKey)
	pk, _ := common.NewPubKey(bech32PubKey)
	return pk
}

// GetRandomTxID create a random tx id for test purpose
func GetRandomTxID() string {
	return common.RandStringBytesMask(64)
}
