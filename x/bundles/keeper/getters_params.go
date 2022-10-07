package keeper

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UploadTimeout(ctx),
		k.StorageCost(ctx),
		k.NetworkFee(ctx),
		k.MaxPoints(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// UploadTimeout returns the UploadTimeout param
func (k Keeper) UploadTimeout(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyUploadTimeout, &res)
	return
}

// StorageCost returns the StorageCost param
func (k Keeper) StorageCost(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyStorageCost, &res)
	return
}

// NetworkFee returns the NetworkFee param
func (k Keeper) NetworkFee(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyNetworkFee, &res)
	return
}

// MaxPoints returns the MaxPoints param
func (k Keeper) MaxPoints(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxPoints, &res)
	return
}

// ParamStore ...
func (k Keeper) ParamStore() (paramStore paramtypes.Subspace) {
	return k.paramstore
}
