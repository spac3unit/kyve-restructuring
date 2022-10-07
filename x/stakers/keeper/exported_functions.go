package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// These functions are meant to be called from external modules
// For now this is the bundles module and the delegation module
// which need to interact stakers module.

// All these functions are safe in the way that they do not return errors
// and every edge case is handled within the function itself.

// GetTotalStake returns the sum of stake of all stakers who are
// currently participating in the given pool
//func (k Keeper) GetTotalStake(ctx sdk.Context, poolId uint64) uint64 {
//	return k.getStat(ctx, poolId, types.STAKER_STATS_TOTAL_STAKE)
//}

// GetAllStakerAddressesOfPool returns a list of all stakers
// which have a currently a valaccount registered for the given pool
func (k Keeper) GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string) {
	for _, valaccount := range k.GetAllValaccountsOfPool(ctx, poolId) {
		stakers = append(stakers, valaccount.Staker)
	}

	return
}

// GetCommission returns the commission of a staker as a parsed sdk.Dec
func (k Keeper) GetCommission(ctx sdk.Context, stakerAddress string) sdk.Dec {
	staker, _ := k.GetStaker(ctx, stakerAddress)
	uploaderCommission, err := sdk.NewDecFromStr(staker.Commission)
	if err != nil {
		util.LogFatalLogicError("Commission not parsable", err.Error(), staker.Commission)
	}
	return uploaderCommission
}

func (k Keeper) GetSlashFraction(ctx sdk.Context, slashType types.SlashType) (slashAmountRatio sdk.Dec) {
	// Retrieve slash fraction from params
	switch slashType {
	case types.SLASH_TYPE_TIMEOUT:
		slashAmountRatio, _ = sdk.NewDecFromStr(k.TimeoutSlash(ctx))
	case types.SLASH_TYPE_VOTE:
		slashAmountRatio, _ = sdk.NewDecFromStr(k.VoteSlash(ctx))
	case types.SLASH_TYPE_UPLOAD:
		slashAmountRatio, _ = sdk.NewDecFromStr(k.UploadSlash(ctx))
	}
	return
}

func (k Keeper) AssertValaccountAuthorized(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) error {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)

	if !found {
		return types.ErrValaccountUnauthorized
	}

	if valaccount.Valaddress != valaddress {
		return types.ErrValaccountUnauthorized
	}

	return nil
}
