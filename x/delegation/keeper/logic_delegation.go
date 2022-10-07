package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegate performs a safe delegation with all necessary checks
// Warning: does not transfer the amount (only the rewards)
func (k Keeper) performDelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string, amount uint64) {

	// Update in-memory staker index for efficient queries
	k.RemoveStakerIndex(ctx, stakerAddress)
	defer k.SetStakerIndex(ctx, stakerAddress)

	if k.DoesDelegatorExist(ctx, stakerAddress, delegatorAddress) {
		// If the sender is already a delegator, first perform an undelegation, before then delegating.
		reward := k.f1WithdrawRewards(ctx, stakerAddress, delegatorAddress)
		err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, delegatorAddress, reward)
		if err != nil {
			util.PanicHalt(k.upgradeKeeper, ctx, "no money left in module")
		}
		// Emit withdraw event.
		ctx.EventManager().EmitTypedEvent(&types.EventWithdrawRewards{
			Address:  delegatorAddress,
			FromNode: stakerAddress,
			Amount:   reward,
		})

		// Perform redelegation
		unDelegateAmount := k.f1RemoveDelegator(ctx, stakerAddress, delegatorAddress)
		k.f1CreateDelegator(ctx, stakerAddress, delegatorAddress, unDelegateAmount+amount)
	} else {
		// If the sender isn't already a delegator, simply create a new delegation entry.
		k.f1CreateDelegator(ctx, stakerAddress, delegatorAddress, amount)
	}
}

// performUndelegation performs immediately an undelegation of the given amount from the given staker
// If the amount is greater than the available amount, only the available amount will be undelegated.
// This method also transfers the rewards back to the given user.
func (k Keeper) performUndelegation(ctx sdk.Context, stakerAddress string, delegatorAddress string, amount uint64) (actualAmount uint64) {

	// Update in-memory staker index for efficient queries
	k.RemoveStakerIndex(ctx, stakerAddress)
	defer k.SetStakerIndex(ctx, stakerAddress)

	// Withdraw all rewards for the sender.
	reward := k.f1WithdrawRewards(ctx, stakerAddress, delegatorAddress)

	// Transfer tokens from this module to sender.
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, delegatorAddress, reward); err != nil {
		util.PanicHalt(k.upgradeKeeper, ctx, "not enough money in module")
	}

	// Emit withdraw event.
	ctx.EventManager().EmitTypedEvent(&types.EventWithdrawRewards{
		Address:  delegatorAddress,
		FromNode: stakerAddress,
		Amount:   reward,
	})

	// Perform an internal re-delegation.
	undelegatedAmount := k.f1RemoveDelegator(ctx, stakerAddress, delegatorAddress)

	redelegation := uint64(0)
	if undelegatedAmount > amount {
		// if user didnt undelegate everything ...
		redelegation = undelegatedAmount - amount
		// ... create a new delegator entry with the remaining amount
		k.f1CreateDelegator(ctx, stakerAddress, delegatorAddress, redelegation)
	}

	return undelegatedAmount - redelegation
}
