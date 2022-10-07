package keeper

import (
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingDelegationTime(ctx),
		k.RedelegationCooldown(ctx),
		k.RedelegationMaxAmount(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// UnbondingDelegationTime ...
func (k Keeper) UnbondingDelegationTime(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyUnbondingDelegationTime, &res)
	return
}

// RedelegationCooldown ...
func (k Keeper) RedelegationCooldown(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRedelegationCooldown, &res)
	return
}

// RedelegationMaxAmount ...
func (k Keeper) RedelegationMaxAmount(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRedelegationMaxAmount, &res)
	return
}

// ParamStore ...
func (k Keeper) ParamStore() (paramStore paramtypes.Subspace) {
	return k.paramstore
}
