package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/bundles/types"
	stakersmoduletypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleUploadTimeout is an end block hook that triggers an upload timeout for every pool (if applicable).
func (k Keeper) HandleUploadTimeout(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Iterate over all pool Ids.
	for _, pool := range k.poolKeeper.GetAllPools(ctx) {
		err := k.AssertPoolCanRun(ctx, pool.Id)
		bundleProposal, _ := k.GetBundleProposal(ctx, pool.Id)

		// Remove next uploader if pool is not active
		if err != nil {
			bundleProposal.NextUploader = ""
			k.SetBundleProposal(ctx, bundleProposal)
			continue
		}

		// Skip if we haven't reached the upload interval.
		if uint64(ctx.BlockTime().Unix()) < (bundleProposal.CreatedAt + pool.UploadInterval) {
			continue
		}

		// Check if bundle needs to be dropped
		if bundleProposal.StorageId != "" {
			// check if the quorum was actually reached
			voteDistribution := k.GetVoteDistribution(ctx, pool.Id)

			if voteDistribution.Status == types.BUNDLE_STATUS_NO_QUORUM {
				// handle stakers who did not vote at all
				k.handleNonVoters(ctx, pool.Id)

				// Get next uploader
				voters := append(bundleProposal.VotersValid, bundleProposal.VotersInvalid...)
				nextUploader := ""

				if len(voters) > 0 {
					nextUploader = k.chooseNextUploaderFromSelectedStakers(ctx, pool.Id, voters)
				} else {
					nextUploader = k.chooseNextUploaderFromAllStakers(ctx, pool.Id)
				}

				// If consensus wasn't reached, we drop the bundle and emit an event.
				k.dropCurrentBundleProposal(ctx, pool, bundleProposal, voteDistribution, nextUploader)
				continue
			}
		}

		// Skip if we haven't reached the upload timeout.
		if uint64(ctx.BlockTime().Unix()) < (bundleProposal.CreatedAt + pool.UploadInterval + k.UploadTimeout(ctx)) {
			continue
		}

		// We now know that the pool is active and the upload timeout has been reached.

		// Now we slash and remove the current next_uploader
		// (if he is still participating in the pool) and select a new one.
		_, foundNextUploader := k.stakerKeeper.GetValaccount(ctx, pool.Id, bundleProposal.NextUploader)

		if foundNextUploader {
			k.delegationKeeper.SlashDelegators(ctx, bundleProposal.NextUploader, stakersmoduletypes.SLASH_TYPE_TIMEOUT)
			k.stakerKeeper.RemoveValaccountFromPool(ctx, pool.Id, bundleProposal.NextUploader)
		}

		// Update bundle proposal
		bundleProposal.NextUploader = k.chooseNextUploaderFromAllStakers(ctx, pool.Id)
		bundleProposal.CreatedAt = uint64(ctx.BlockTime().Unix())

		k.SetBundleProposal(ctx, bundleProposal)
	}
}
