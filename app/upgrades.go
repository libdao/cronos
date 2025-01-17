package app

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

func (app *App) RegisterUpgradeHandlers() {
	upgradeHandlerV2 := func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		m, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
		if err != nil {
			return m, err
		}

		gravParams := app.GravityKeeper.GetParams(ctx)
		gravParams.GravityId = "cronos_gravity_testnet"
		// can be activated later on
		gravParams.BridgeActive = false
		app.GravityKeeper.SetParams(ctx, gravParams)
		return m, nil
	}
	// `v1.0.0` upgrade plan will clear the `extra_eips` parameters, and upgrade ibc-go to v5.2.0.
	planName := "v2.0.0-testnet3"
	app.UpgradeKeeper.SetUpgradeHandler(planName, upgradeHandlerV2)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		if upgradeInfo.Name == planName {
			storeUpgrades := storetypes.StoreUpgrades{
				Added: []string{gravitytypes.StoreKey},
			}

			// configure store loader that checks if version == upgradeHeight and applies store upgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
		}
	}
}
