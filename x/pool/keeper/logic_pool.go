package keeper

import (
	"github.com/KYVENetwork/chain/x/pool/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetPoolWithError ...
func (k Keeper) GetPoolWithError(ctx sdk.Context, poolId uint64) (types.Pool, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.Pool{}, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
	}
	return pool, nil
}

// AssertPoolExists ...
func (k Keeper) AssertPoolExists(ctx sdk.Context, poolId uint64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolKey)
	if store.Has(types.PoolKeyPrefix(poolId)) {
		return nil
	}
	return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), poolId)
}

func (k Keeper) IncrementBundleInformation(
	ctx sdk.Context,
	poolId uint64,
	currentHeight uint64,
	currentKey string,
	currentValue string,
) {
	pool, found := k.GetPool(ctx, poolId)
	if found {
		pool.CurrentHeight = currentHeight
		pool.TotalBundles = pool.TotalBundles + 1
		pool.CurrentKey = currentKey
		pool.CurrentValue = currentValue
		k.SetPool(ctx, pool)
	}

}
