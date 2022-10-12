package keeper

import (
	"mercury/x/mercury/types"
)

var _ types.QueryServer = Keeper{}
