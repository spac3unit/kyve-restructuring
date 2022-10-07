package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/x/bundles/types"
)

// ClaimUploaderRole handles the logic of an SDK message that allows protocol nodes to claim the uploader role.
// Note that this function can only be called when the next uploader is not chosen yet.
// This function obeys "first come, first serve" mentality.
func (k msgServer) ClaimUploaderRole(
	goCtx context.Context, msg *types.MsgClaimUploaderRole,
) (*types.MsgClaimUploaderRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if poolErr := k.AssertPoolCanRun(ctx, msg.PoolId); poolErr != nil {
		return nil, poolErr
	}

	if err := k.stakerKeeper.AssertValaccountAuthorized(ctx, msg.PoolId, msg.Staker, msg.Creator); err != nil {
		return nil, err
	}

	bundleProposal, found := k.GetBundleProposal(ctx, msg.PoolId)

	// If the pool was newly created no bundle proposal exists yet.
	// There is one bundle proposal per pool.
	if !found {
		bundleProposal.PoolId = msg.PoolId
	}

	// Error if the next uploader is already set.
	if bundleProposal.NextUploader != "" {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrUploaderAlreadyClaimed.Error())
	}

	bundleProposal.NextUploader = msg.Staker
	bundleProposal.CreatedAt = uint64(ctx.BlockTime().Unix())

	k.SetBundleProposal(ctx, bundleProposal)

	return &types.MsgClaimUploaderRoleResponse{}, nil
}
