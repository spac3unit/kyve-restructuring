package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitBundleProposal handles the logic of an SDK message that allows protocol nodes to submit a new bundle proposal.
func (k msgServer) SubmitBundleProposal(
	goCtx context.Context, msg *types.MsgSubmitBundleProposal,
) (*types.MsgSubmitBundleProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.AssertCanPropose(ctx, msg.PoolId, msg.Staker, msg.Creator, msg.FromHeight); err != nil {
		return nil, err
	}

	pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)

	// Validate submit bundle args.
	if err := k.validateSubmitBundleArgs(ctx, &bundleProposal, msg); err != nil {
		return nil, err
	}

	// reset points of uploader as node has proven to be active
	k.stakerKeeper.ResetPoints(ctx, msg.PoolId, msg.Staker)

	// If bundle was dropped just register the new bundle.
	if bundleProposal.StorageId == "" {
		nextUploader := k.chooseNextUploaderFromAllStakers(ctx, msg.PoolId)

		if err := k.registerBundleProposalFromUploader(ctx, pool, bundleProposal, msg, nextUploader); err != nil {
			return nil, err
		}

		return &types.MsgSubmitBundleProposalResponse{}, nil
	}

	// Previous round contains a bundle which needs to be validated now.

	// increase points of stakers who did not vote at all + slash + remove if necessary
	k.handleNonVoters(ctx, msg.PoolId)

	// Get next uploader from stakers voted
	voters := append(bundleProposal.VotersValid, bundleProposal.VotersInvalid...)
	nextUploader := ""

	if len(voters) > 0 {
		nextUploader = k.chooseNextUploaderFromSelectedStakers(ctx, msg.PoolId, voters)
	} else {
		nextUploader = k.chooseNextUploaderFromAllStakers(ctx, msg.PoolId)
	}

	// evaluate all votes and determine status based on the votes weighted with stake + delegation
	voteDistribution := k.GetVoteDistribution(ctx, msg.PoolId)

	// handle valid proposal
	if voteDistribution.Status == types.BUNDLE_STATUS_VALID {
		// Calculate the total reward for the bundle, and individual payouts.
		bundleReward := k.calculatePayouts(ctx, msg.PoolId)

		if err := k.poolKeeper.ChargeFundersOfPool(ctx, msg.PoolId, bundleReward.Total); err != nil {
			// drop bundle because pool ran out of funds
			bundleProposal.CreatedAt = uint64(ctx.BlockTime().Unix())
			k.SetBundleProposal(ctx, bundleProposal)

			// emit event which indicates that pool has run out of funds
			if errEmit := ctx.EventManager().EmitTypedEvent(&pooltypes.EventPoolOutOfFunds{
				PoolId: msg.PoolId,
			}); errEmit != nil {
				return nil, errEmit
			}

			return &types.MsgSubmitBundleProposalResponse{}, nil
		}

		pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)
		bundleProposal, _ := k.GetBundleProposal(ctx, msg.PoolId)

		uploaderPayout := bundleReward.Uploader

		delegationPayoutSuccessful := k.delegationKeeper.PayoutRewards(ctx, bundleProposal.Uploader, bundleReward.Delegation, pooltypes.ModuleName)
		// If staker has no delegators add all delegation rewards to the staker rewards
		if !delegationPayoutSuccessful {
			uploaderPayout += bundleReward.Delegation
		}

		// send commission to uploader
		if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, pooltypes.ModuleName, bundleProposal.Uploader, uploaderPayout); err != nil {
			return nil, err
		}

		// send network fee to treasury
		if err := util.TransferFromModuleToTreasury(k.accountKeeper, k.distrkeeper, ctx, pooltypes.ModuleName, bundleReward.Treasury); err != nil {
			return nil, err
		}

		// slash stakers who voted incorrectly
		for _, voter := range bundleProposal.VotersInvalid {
			k.delegationKeeper.SlashDelegators(ctx, voter, stakertypes.SLASH_TYPE_VOTE)
		}

		if err := k.finalizeCurrentBundleProposal(ctx, pool, bundleProposal, voteDistribution, bundleReward); err != nil {
			return nil, err
		}

		// Register the provided bundle as a new proposal for the next round
		if err := k.registerBundleProposalFromUploader(ctx, pool, bundleProposal, msg, nextUploader); err != nil {
			return nil, err
		}

		return &types.MsgSubmitBundleProposalResponse{}, nil
	} else if voteDistribution.Status == types.BUNDLE_STATUS_INVALID {
		// slash stakers who voted incorrectly - uploader receives upload slash
		for _, voter := range bundleProposal.VotersValid {
			if voter == bundleProposal.Uploader {
				k.delegationKeeper.SlashDelegators(ctx, voter, stakertypes.SLASH_TYPE_UPLOAD)
			} else {
				k.delegationKeeper.SlashDelegators(ctx, voter, stakertypes.SLASH_TYPE_VOTE)
			}
		}

		// Drop current bundle. Can't register the provided bundle because the previous bundles
		// needs to be resubmitted first.
		if err := k.dropCurrentBundleProposal(ctx, pool, bundleProposal, voteDistribution, bundleProposal.NextUploader); err != nil {
			return nil, err
		}

		return &types.MsgSubmitBundleProposalResponse{}, nil
	} else {
		return nil, types.ErrQuorumNotReached
	}
}
