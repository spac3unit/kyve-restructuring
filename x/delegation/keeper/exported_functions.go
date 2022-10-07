package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// These functions are meant to be called from external modules
// For now this is the bundles module which needs to interact
// with the delegation module.
// All these functions are safe in the way that they do not return errors
// and every edge case is handled within the function itself.

// GetDelegationAmount returns the sum of all delegations for a specific staker.
// If the staker does not exist, it returns zero as the staker as zero delegations
func (k Keeper) GetDelegationAmount(ctx sdk.Context, staker string) uint64 {
	delegationData, found := k.GetDelegationData(ctx, staker)

	if found {
		return delegationData.TotalDelegation
	}

	return 0
}

// GetDelegationAmountOfDelegator returns the amount of how many $KYVE `delegatorAddress`
// has delegated to `stakerAddress`. If one of the addresses does not exist, it returns zero.
func (k Keeper) GetDelegationAmountOfDelegator(ctx sdk.Context, stakerAddress string, delegatorAddress string) uint64 {
	return k.f1GetCurrentDelegation(ctx, stakerAddress, delegatorAddress)
}

// GetDelegationOfPool returns the amount of how many $KYVE users have delegated
// to stakers that are participating in the given pool
func (k Keeper) GetDelegationOfPool(ctx sdk.Context, poolId uint64) uint64 {
	totalDelegation := uint64(0)
	for _, address := range k.stakersKeeper.GetAllStakerAddressesOfPool(ctx, poolId) {
		totalDelegation += k.GetDelegationAmount(ctx, address)
	}
	return totalDelegation
}

// PayoutRewards transfers `amount` $nKYVE from the `payerModuleName`-module to the delegation module.
// It then awards these tokens internally to all delegators of staker `staker`.
// Delegators can then receive these rewards if they call the `withdraw`-transaction.
// This method return false if the payout failed. This happens usually if there are no
// delegators for that staker. If this happens one should do something else with the rewards.
func (k Keeper) PayoutRewards(ctx sdk.Context, staker string, amount uint64, payerModuleName string) (success bool) {
	// Assert there are delegators
	if k.DoesDelegationDataExist(ctx, staker) {

		// Add amount to the rewards pool
		k.AddAmountToDelegationRewards(ctx, staker, amount)

		// Transfer tokens to the delegation module
		err := util.TransferFromModuleToModule(k.bankKeeper, ctx, payerModuleName, types.ModuleName, amount)
		if err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "Not enough tokens in module")
			return false
		}
		return true
	}
	return false
}

// SlashDelegators reduces the delegation of all delegators of `staker` by fraction
// and transfers the amount to the Treasury.
func (k Keeper) SlashDelegators(ctx sdk.Context, staker string, slashType stakerstypes.SlashType) {

	// Only slash if staker has delegators
	if k.DoesDelegationDataExist(ctx, staker) {

		// Update in-memory staker index for efficient queries
		k.RemoveStakerIndex(ctx, staker)
		defer k.SetStakerIndex(ctx, staker)

		// Perform F1-slash and get slashed amount in nKYVE
		slashedAmount := k.f1Slash(ctx, staker, k.stakersKeeper.GetSlashFraction(ctx, slashType))

		// Transfer tokens to the Treasury
		if err := util.TransferFromModuleToTreasury(k.accountKeeper, k.distrKeeper, ctx, types.ModuleName, slashedAmount); err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "Not enough tokens in module")
		}
	}

}

// GetOutstandingRewards calculates the current rewards a delegator has collected for
// the given staker.
func (k Keeper) GetOutstandingRewards(ctx sdk.Context, staker string, delegator string) uint64 {
	return k.f1GetOutstandingRewards(ctx, staker, delegator)
}
