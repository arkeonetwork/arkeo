package keeper

import (
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

var _ types.QueryServer = Keeper{}
