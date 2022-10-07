package keeper

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math"
)

func (k Keeper) getLowestStaker(ctx sdk.Context, poolId uint64) (val types.Staker, found bool) {
	var minAmount uint64 = math.MaxUint64

	for _, staker := range k.getAllStakersOfPool(ctx, poolId) {
		delegationAmount := k.delegationKeeper.GetDelegationAmount(ctx, staker.Address)
		if delegationAmount < minAmount {
			minAmount = delegationAmount
			val = staker
		}
	}

	return
}

func (k Keeper) ensureFreeSlot(ctx sdk.Context, poolId uint64, stakerAddress string) error {
	// check if slots are still available
	if k.GetStakerCountOfPool(ctx, poolId) >= types.MaxStakers {
		// if not - get lowest staker
		lowestStaker, _ := k.getLowestStaker(ctx, poolId)

		// if new pool joiner has more stake than lowest staker kick him out
		if k.delegationKeeper.GetDelegationAmount(ctx, stakerAddress) > k.delegationKeeper.GetDelegationAmount(ctx, lowestStaker.Address) {
			// remove lowest staker from pool
			k.RemoveValaccountFromPool(ctx, poolId, lowestStaker.Address)

			// emit event
			if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventLeavePool{
				PoolId: poolId,
				Staker: lowestStaker.Address,
			}); errEmit != nil {
				return errEmit
			}
		} else {
			return sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrStakeTooLow.Error(), k.delegationKeeper.GetDelegationAmount(ctx, lowestStaker.Address))
		}
	}

	return nil
}
