package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type AccountKeeper interface {
}

type UpgradeKeeper interface {
	ScheduleUpgrade(ctx sdk.Context, plan upgradetypes.Plan) error
}
