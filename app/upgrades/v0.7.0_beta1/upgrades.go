package v0_7_0_beta1

import (
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(registryKeeper *registrykeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		_ = registryKeeper

		return vm, nil
	}
}
