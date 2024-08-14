package app

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/arkeonetwork/arkeo/app/upgrades"
)

// Upgrades
var Upgrades = []upgrades.Upgrade{}

func (app *ArkeoApp) RegisterUpgradeHandlers() {
	app.setUpgradeHandlers()
	app.setUpgradeStoreLoaders()

}

func (app *ArkeoApp) setUpgradeStoreLoaders() {
	upgradeInfo, err := app.Keepers.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.Keepers.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, u := range Upgrades {
		if upgradeInfo.Name == u.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &u.StoreUpgrades))
		}
	}
}

func (app *ArkeoApp) setUpgradeHandlers() {
	for _, u := range Upgrades {
		app.Keepers.UpgradeKeeper.SetUpgradeHandler(
			u.UpgradeName,
			u.CreateUpgradeHandler(app.mm, app.configurator, app.Keepers),
		)
	}
}
