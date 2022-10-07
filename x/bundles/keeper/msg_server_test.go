package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/KYVENetwork/chain/testutil/keeper"
	"github.com/KYVENetwork/chain/x/bundles/keeper"
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.BundlesKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
