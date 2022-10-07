package keeper

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.VoteSlash(ctx),
		k.UploadSlash(ctx),
		k.TimeoutSlash(ctx),
		k.UnbondingStakingTime(ctx),
		k.CommissionChangeTime(ctx),
		k.LeavePoolTime(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// VoteSlash returns the VoteSlash param
func (k Keeper) VoteSlash(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyVoteSlash, &res)
	return
}

// UploadSlash returns the UploadSlash param
func (k Keeper) UploadSlash(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyUploadSlash, &res)
	return
}

// TimeoutSlash returns the TimeoutSlash param
func (k Keeper) TimeoutSlash(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyTimeoutSlash, &res)
	return
}

// UnbondingStakingTime returns the UnbondingStakingTime param
func (k Keeper) UnbondingStakingTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyUnbondingStakingTime, &res)
	return
}

// CommissionChangeTime returns the CommissionChangeTime param
func (k Keeper) CommissionChangeTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyCommissionChangeTime, &res)
	return
}

// LeavePoolTime returns the LeavePoolTime param
func (k Keeper) LeavePoolTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyLeavePoolTime, &res)
	return
}

// ParamStore returns the entire param store
func (k Keeper) ParamStore() (paramStore paramtypes.Subspace) {
	return k.paramstore
}
