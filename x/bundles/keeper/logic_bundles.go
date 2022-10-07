package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/bundles/types"
	poolmoduletypes "github.com/KYVENetwork/chain/x/pool/types"
	stakermoduletypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math"
	"math/rand"
	"sort"
)

// AssertPoolCanRun checks whether the given pool fulfils all
// technical/formal requirements to produce bundles
func (k Keeper) AssertPoolCanRun(ctx sdk.Context, poolId uint64) error {

	pool, poolErr := k.poolKeeper.GetPoolWithError(ctx, poolId)
	if poolErr != nil {
		return poolErr
	}

	// Error if the pool is upgrading.
	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		return types.ErrPoolCurrentlyUpgrading
	}

	// Error if the pool is paused.
	if pool.Paused {
		return types.ErrPoolPaused
	}

	// Error if the pool has no funds.
	if len(pool.Funders) == 0 {
		return types.ErrPoolOutOfFunds
	}

	// Error if min stake is not reached
	if k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id) < pool.MinStake {
		return types.ErrMinStakeNotReached
	}

	return nil
}

func (k Keeper) AssertCanVote(ctx sdk.Context, poolId uint64, staker string, voter string, storageId string) error {
	// Check basic pool configs
	if err := k.AssertPoolCanRun(ctx, poolId); err != nil {
		return err
	}

	// Check if sender is a staker in pool
	if err := k.stakerKeeper.AssertValaccountAuthorized(ctx, poolId, staker, voter); err != nil {
		return err
	}

	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Check if dropped bundle
	if bundleProposal.StorageId == "" {
		return types.ErrBundleDropped
	}

	// Check if tx matches current bundleProposal
	if storageId != bundleProposal.StorageId {
		return types.ErrInvalidStorageId
	}

	// Check if the sender has already voted on the bundle.
	hasVotedValid := util.ContainsString(bundleProposal.VotersValid, staker)
	hasVotedInvalid := util.ContainsString(bundleProposal.VotersInvalid, staker)

	if hasVotedValid {
		return types.ErrAlreadyVotedValid
	}

	if hasVotedInvalid {
		return types.ErrAlreadyVotedInvalid
	}

	return nil
}

func (k Keeper) AssertCanPropose(ctx sdk.Context, poolId uint64, staker string, proposer string, fromHeight uint64) error {
	// Check basic pool configs
	if err := k.AssertPoolCanRun(ctx, poolId); err != nil {
		return err
	}

	// Check if sender is a staker in pool
	if err := k.stakerKeeper.AssertValaccountAuthorized(ctx, poolId, staker, proposer); err != nil {
		return err
	}

	pool, _ := k.poolKeeper.GetPool(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Check if designated uploader
	if bundleProposal.NextUploader != staker {
		return sdkErrors.Wrapf(types.ErrNotDesignatedUploader, "expected %v received %v", bundleProposal.NextUploader, staker)
	}

	// Check if upload interval has been surpassed
	if uint64(ctx.BlockTime().Unix()) < (bundleProposal.CreatedAt + pool.UploadInterval) {
		return sdkErrors.Wrapf(types.ErrUploadInterval, "expected %v < %v", ctx.BlockTime().Unix(), bundleProposal.CreatedAt+pool.UploadInterval)
	}

	// Check if from_height matches
	if bundleProposal.ToHeight == 0 {
		if pool.CurrentHeight != fromHeight {
			return sdkErrors.Wrapf(types.ErrFromHeight, "expected %v received %v", pool.CurrentHeight, fromHeight)
		}
	} else {
		if bundleProposal.ToHeight != fromHeight {
			return sdkErrors.Wrapf(types.ErrFromHeight, "expected %v received %v", bundleProposal.ToHeight, fromHeight)
		}
	}

	return nil
}

func (k Keeper) validateSubmitBundleArgs(ctx sdk.Context, bundleProposal *types.BundleProposal, msg *types.MsgSubmitBundleProposal) error {
	pool, _ := k.poolKeeper.GetPool(ctx, msg.PoolId)

	currentHeight := bundleProposal.ToHeight

	if currentHeight == 0 {
		currentHeight = pool.CurrentHeight
	}

	currentKey := bundleProposal.ToKey

	// Validate storage id
	if msg.StorageId == "" {
		return types.ErrInvalidArgs
	}

	// Validate if bundle is not too big
	if msg.ToHeight-currentHeight > pool.MaxBundleSize {
		return sdkErrors.Wrapf(types.ErrMaxBundleSize, "expected %v received %v", pool.MaxBundleSize, msg.ToHeight-currentHeight)
	}

	// Validate to height
	if msg.ToHeight < currentHeight {
		return sdkErrors.Wrapf(types.ErrToHeight, "expected %v < %v", msg.ToHeight, currentHeight)
	}

	// Validate from key
	if currentKey != "" && msg.FromKey != currentKey {
		return types.ErrFromKey
	}

	// Validate height and byte size
	if msg.ToHeight <= currentHeight || msg.ByteSize == 0 {
		return types.ErrInvalidArgs
	}

	// Validate key values
	if msg.ToKey == "" || msg.ToValue == "" {
		return types.ErrInvalidArgs
	}

	// Validate bundle hash
	if msg.BundleHash == "" {
		return types.ErrInvalidArgs
	}

	return nil
}

func (k Keeper) registerBundleProposalFromUploader(ctx sdk.Context, pool poolmoduletypes.Pool, bundleProposal types.BundleProposal, msg *types.MsgSubmitBundleProposal, nextUploader string) error {
	bundleProposal = types.BundleProposal{
		PoolId:       msg.PoolId,
		Uploader:     msg.Staker,
		NextUploader: nextUploader,
		StorageId:    msg.StorageId,
		ByteSize:     msg.ByteSize,
		ToHeight:     msg.ToHeight,
		CreatedAt:    uint64(ctx.BlockTime().Unix()),
		VotersValid:  append(make([]string, 0), msg.Staker),
		ToKey:        msg.ToKey,
		ToValue:      msg.ToValue,
		BundleHash:   msg.BundleHash,
	}

	k.SetBundleProposal(ctx, bundleProposal)

	err := ctx.EventManager().EmitTypedEvent(&types.EventBundleProposed{
		PoolId:     bundleProposal.PoolId,
		Id:         pool.TotalBundles,
		StorageId:  bundleProposal.StorageId,
		Uploader:   bundleProposal.Uploader,
		ByteSize:   bundleProposal.ByteSize,
		FromHeight: pool.CurrentHeight,
		ToHeight:   bundleProposal.ToHeight,
		FromKey:    pool.CurrentKey,
		ToKey:      bundleProposal.ToKey,
		Value:      bundleProposal.ToValue,
		BundleHash: bundleProposal.BundleHash,
		CreatedAt:  bundleProposal.CreatedAt,
	})

	return err
}

// handleNonVoters checks if stakers in a pool voted on the current bundle proposal
// if a staker did not vote at all on a bundle proposal he received points
// if a staker receives a certain number of points he receives a timeout slash and gets
// kicked out of a pool
func (k Keeper) handleNonVoters(ctx sdk.Context, poolId uint64) {
	voters := map[string]bool{}
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	for _, address := range bundleProposal.VotersValid {
		voters[address] = true
	}

	for _, address := range bundleProposal.VotersInvalid {
		voters[address] = true
	}

	for _, address := range bundleProposal.VotersAbstain {
		voters[address] = true
	}

	for _, staker := range k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		if !voters[staker] {
			points := k.stakerKeeper.AddPoint(ctx, poolId, staker)

			if points >= k.MaxPoints(ctx) {
				k.delegationKeeper.SlashDelegators(ctx, staker, stakermoduletypes.SLASH_TYPE_TIMEOUT)
				k.stakerKeeper.ResetPoints(ctx, poolId, staker)
				k.stakerKeeper.RemoveValaccountFromPool(ctx, poolId, staker)
			}
		}
	}
}

// calculatePayouts deducts the network fee from the rewards and splits the remaining amount
// between the staker and its delegators. If there are no delegators, the entire amount is
// awarded to the staker.
func (k Keeper) calculatePayouts(ctx sdk.Context, poolId uint64) (bundleReward types.BundleReward) {

	pool, _ := k.poolKeeper.GetPool(ctx, poolId)
	bundleProposal, _ := k.GetBundleProposal(ctx, poolId)

	// Should not happen, if so move everything to the treasury
	if !k.stakerKeeper.DoesStakerExist(ctx, bundleProposal.Uploader) {
		bundleReward.Treasury = bundleReward.Total

		return
	}

	// formula for calculating the rewards
	bundleReward.Total = pool.OperatingCost + (bundleProposal.ByteSize * k.StorageCost(ctx))

	networkFee, err := sdk.NewDecFromStr(k.NetworkFee(ctx))
	if err != nil {
		util.LogFatalLogicError("Network Fee unparasable", err.Error(), k.NetworkFee(ctx))
	}
	// Add fee to treasury
	bundleReward.Treasury = uint64(sdk.NewDec(int64(bundleReward.Total)).Mul(networkFee).RoundInt64())

	// Remaining rewards to be split between staker and its delegators
	totalNodeReward := bundleReward.Total - bundleReward.Treasury

	// Payout delegators
	if k.delegationKeeper.GetDelegationAmount(ctx, bundleProposal.Uploader) > 0 {
		commission := k.stakerKeeper.GetCommission(ctx, bundleProposal.Uploader)

		bundleReward.Uploader = uint64(sdk.NewDec(int64(totalNodeReward)).Mul(commission).RoundInt64())
		bundleReward.Delegation = totalNodeReward - bundleReward.Uploader
	} else {
		bundleReward.Uploader = totalNodeReward
		bundleReward.Delegation = 0
	}

	return
}

func (k Keeper) finalizeCurrentBundleProposal(ctx sdk.Context, pool poolmoduletypes.Pool, bundleProposal types.BundleProposal, voteDistribution types.VoteDistribution, bundleReward types.BundleReward) error {
	// save finalized bundle
	finalizedBundle := types.FinalizedBundle{
		StorageId:   bundleProposal.StorageId,
		PoolId:      pool.Id,
		Id:          pool.TotalBundles,
		Uploader:    bundleProposal.Uploader,
		FromHeight:  pool.CurrentHeight,
		ToHeight:    bundleProposal.ToHeight,
		FinalizedAt: uint64(ctx.BlockHeight()),
		Key:         bundleProposal.ToKey,
		Value:       bundleProposal.ToValue,
		BundleHash:  bundleProposal.BundleHash,
	}

	k.SetFinalizedBundle(ctx, finalizedBundle)

	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalized{
		PoolId:           finalizedBundle.PoolId,
		Id:               finalizedBundle.Id,
		Valid:            voteDistribution.Valid,
		Invalid:          voteDistribution.Invalid,
		Abstain:          voteDistribution.Abstain,
		Total:            voteDistribution.Total,
		Status:           voteDistribution.Status,
		RewardTreasury:   bundleReward.Treasury,
		RewardUploader:   bundleReward.Uploader,
		RewardDelegation: bundleReward.Delegation,
		RewardTotal:      bundleReward.Total,
	}); errEmit != nil {
		return errEmit
	}

	// Finalize the proposal, saving useful information.
	k.poolKeeper.IncrementBundleInformation(ctx, pool.Id, bundleProposal.ToHeight, bundleProposal.ToKey, bundleProposal.ToValue)

	return nil
}

func (k Keeper) dropCurrentBundleProposal(
	ctx sdk.Context, pool poolmoduletypes.Pool,
	bundleProposal types.BundleProposal,
	voteDistribution types.VoteDistribution,
	nextUploader string,
) error {
	err := ctx.EventManager().EmitTypedEvent(&types.EventBundleFinalized{
		PoolId:  pool.Id,
		Id:      pool.TotalBundles,
		Valid:   voteDistribution.Valid,
		Invalid: voteDistribution.Invalid,
		Abstain: voteDistribution.Abstain,
		Total:   voteDistribution.Total,
		Status:  voteDistribution.Status,
	})

	// drop bundle
	bundleProposal = types.BundleProposal{
		PoolId:       pool.Id,
		NextUploader: nextUploader,
		CreatedAt:    uint64(ctx.BlockTime().Unix()),
	}

	k.SetBundleProposal(ctx, bundleProposal)

	return err
}

// RandomChoiceCandidate ...
type RandomChoiceCandidate struct {
	Account string
	Amount  uint64
}

// getWeightedRandomChoice is an internal function that returns a random selection out of a list of candidates.
func (k Keeper) getWeightedRandomChoice(candidates []RandomChoiceCandidate, seed uint64) string {
	type WeightedRandomChoice struct {
		Elements    []string
		Weights     []uint64
		TotalWeight uint64
	}

	wrc := WeightedRandomChoice{}

	for _, candidate := range candidates {
		i := sort.Search(len(wrc.Weights), func(i int) bool { return wrc.Weights[i] > candidate.Amount })
		wrc.Weights = append(wrc.Weights, 0)
		wrc.Elements = append(wrc.Elements, "")
		copy(wrc.Weights[i+1:], wrc.Weights[i:])
		copy(wrc.Elements[i+1:], wrc.Elements[i:])
		wrc.Weights[i] = candidate.Amount
		wrc.Elements[i] = candidate.Account
		wrc.TotalWeight += candidate.Amount
	}

	rand.Seed(int64(seed))
	value := uint64(math.Floor(rand.Float64() * float64(wrc.TotalWeight)))

	for key, weight := range wrc.Weights {
		if weight > value {
			return wrc.Elements[key]
		}

		value -= weight
	}

	return ""
}

func (k Keeper) chooseNextUploaderFromSelectedStakers(ctx sdk.Context, poolId uint64, addresses []string) (nextUploader string) {
	var _candidates []RandomChoiceCandidate

	if len(addresses) == 0 {
		return ""
	}

	for _, s := range addresses {
		delegation := k.delegationKeeper.GetDelegationAmount(ctx, s)

		_candidates = append(_candidates, RandomChoiceCandidate{
			Account: s,
			Amount:  delegation,
		})
	}

	return k.getWeightedRandomChoice(_candidates, uint64(ctx.BlockHeight()))
}

func (k Keeper) chooseNextUploaderFromAllStakers(ctx sdk.Context, poolId uint64) (nextUploader string) {
	stakers := k.stakerKeeper.GetAllStakerAddressesOfPool(ctx, poolId)
	return k.chooseNextUploaderFromSelectedStakers(ctx, poolId, stakers)
}

// getVoteDistribution is an internal function evaulates the quorum status of a bundle proposal.
func (k Keeper) GetVoteDistribution(ctx sdk.Context, poolId uint64) (voteDistribution types.VoteDistribution) {
	bundleProposal, found := k.GetBundleProposal(ctx, poolId)
	if !found {
		return
	}

	// get $KYVE voted for valid
	for _, voter := range bundleProposal.VotersValid {
		delegation := k.delegationKeeper.GetDelegationAmount(ctx, voter)
		voteDistribution.Valid += delegation
	}

	// get $KYVE voted for invalid
	for _, voter := range bundleProposal.VotersInvalid {
		delegation := k.delegationKeeper.GetDelegationAmount(ctx, voter)
		voteDistribution.Invalid += delegation
	}

	// get $KYVE voted for abstain
	for _, voter := range bundleProposal.VotersAbstain {
		delegation := k.delegationKeeper.GetDelegationAmount(ctx, voter)
		voteDistribution.Abstain += delegation
	}

	voteDistribution.Total = k.delegationKeeper.GetDelegationOfPool(ctx, poolId)

	if voteDistribution.Valid*2 > voteDistribution.Total {
		voteDistribution.Status = types.BUNDLE_STATUS_VALID
	} else if voteDistribution.Invalid*2 >= voteDistribution.Total {
		voteDistribution.Status = types.BUNDLE_STATUS_INVALID
	} else {
		voteDistribution.Status = types.BUNDLE_STATUS_NO_QUORUM
	}

	return
}
