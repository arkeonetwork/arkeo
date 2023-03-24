package types

import (
	"github.com/arkeonetwork/arkeo/common"
	arkeoTypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

// wrapper for the arkeo contract so that we can included nonces with the contract
// for our in memory store.
type SentinelContract struct {
	ArkeoContract arkeoTypes.Contract
	Nonces        map[common.PubKey]int64
}
