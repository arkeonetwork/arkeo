package upgrade1_0_0

import (
	"context"

	storetypes "cosmossdk.io/store/types"

	"cosmossdk.io/x/upgrade/types"
	"github.com/arkeonetwork/arkeo/app/upgrades"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const Name = "v1.0.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName: Name,
	CreateUpgradeHandler: func(m *module.Manager, cfg module.Configurator) types.UpgradeHandler {
		return func(ctx context.Context, plan types.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

			return m.RunMigrations(ctx, cfg, fromVM)
		}
	},
	StoreUpgrades: storetypes.StoreUpgrades{},
}
