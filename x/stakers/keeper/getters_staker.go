package keeper

import (
	"encoding/binary"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateStakerMetadata ...
func (k Keeper) UpdateStakerMetadata(ctx sdk.Context, address string, moniker string, website string, logo string) {
	staker, found := k.GetStaker(ctx, address)
	if found {
		staker.Moniker = moniker
		staker.Website = website
		staker.Logo = logo
		k.setStaker(ctx, staker)
	}
}

// UpdateStakerCommission ...
func (k Keeper) UpdateStakerCommission(ctx sdk.Context, address string, commission string) {
	staker, found := k.GetStaker(ctx, address)
	if found {
		staker.Commission = commission
		k.setStaker(ctx, staker)
	}
}

// AddValaccountToPool adds a valaccount to a pool.
// If valaccount already belongs to pool, nothing happens.
func (k Keeper) AddValaccountToPool(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) {
	if k.DoesStakerExist(ctx, stakerAddress) {
		if !k.DoesValaccountExist(ctx, poolId, stakerAddress) {
			k.SetValaccount(ctx, types.Valaccount{
				PoolId:     poolId,
				Staker:     stakerAddress,
				Valaddress: valaddress,
			})
			k.AddOneToCount(ctx, poolId)
		}
	}
}

// RemoveValaccountFromPool removes a valaccount from a given pool and updates
// all aggregated variables. If the valaccount is not in the pool nothing happens.
func (k Keeper) RemoveValaccountFromPool(ctx sdk.Context, poolId uint64, stakerAddress string) {
	// get valaccount
	valaccount, valaccountFound := k.GetValaccount(ctx, poolId, stakerAddress)

	// if valaccount was found on pool continue
	if valaccountFound {
		// remove valaccount from pool
		k.removeValaccount(ctx, valaccount)
		k.subtractOneFromCount(ctx, poolId)
	}
}

func (k Keeper) AppendStaker(ctx sdk.Context, staker types.Staker) {
	k.setStaker(ctx, staker)
}

// #############################
// #  Raw KV-Store operations  #
// #############################

func (k Keeper) getAllStakersOfPool(ctx sdk.Context, poolId uint64) []types.Staker {
	valaccounts := k.GetAllValaccountsOfPool(ctx, poolId)

	stakers := make([]types.Staker, 0)

	for _, valaccount := range valaccounts {
		staker, _ := k.GetStaker(ctx, valaccount.Staker)
		stakers = append(stakers, staker)
	}

	return stakers
}

// removeStaker removes a staker from the store
func (k Keeper) removeStaker(ctx sdk.Context, staker types.Staker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StakerKeyPrefix)
	store.Delete(types.StakerKey(
		staker.Address,
	))
}

// SetStaker set a specific staker in the store from its index
func (k Keeper) setStaker(ctx sdk.Context, staker types.Staker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StakerKeyPrefix)
	b := k.cdc.MustMarshal(&staker)
	store.Set(types.StakerKey(
		staker.Address,
	), b)
}

// GetStaker returns a staker from its index
func (k Keeper) GetStaker(
	ctx sdk.Context,
	staker string,
) (val types.Staker, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StakerKeyPrefix)

	b := store.Get(types.StakerKey(
		staker,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetPaginatedStakerQuery(ctx sdk.Context, pagination *query.PageRequest, accumulator func(staker types.Staker)) (*query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StakerKeyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var staker types.Staker
			if err := k.cdc.Unmarshal(value, &staker); err != nil {
				return false, err
			}
			accumulator(staker)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}

// DoesStakerExist returns true if the staker exists
func (k Keeper) DoesStakerExist(ctx sdk.Context, staker string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StakerKeyPrefix)
	return store.Has(types.StakerKey(staker))
}

// GetAllStakers returns all staker
func (k Keeper) GetAllStakers(ctx sdk.Context) (list []types.Staker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StakerKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Staker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// #############################
// #     Aggregation Data      #
// #############################

func (k Keeper) GetStakerCountOfPool(ctx sdk.Context, poolId uint64) uint64 {
	return k.getStat(ctx, poolId, types.STAKER_STATS_COUNT)
}

func (k Keeper) AddOneToCount(ctx sdk.Context, poolId uint64) {
	count := k.getStat(ctx, poolId, types.STAKER_STATS_COUNT)
	k.setStat(ctx, poolId, types.STAKER_STATS_COUNT, count+1)
}

func (k Keeper) subtractOneFromCount(ctx sdk.Context, poolId uint64) {
	count := k.getStat(ctx, poolId, types.STAKER_STATS_COUNT)
	k.setStat(ctx, poolId, types.STAKER_STATS_COUNT, count-1)
}

// getStat get the total number of pool
func (k Keeper) getStat(ctx sdk.Context, poolId uint64, statType types.STAKER_STATS) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := store.Get(util.GetByteKey(string(statType), poolId))
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// setStat set the total number of pool
func (k Keeper) setStat(ctx sdk.Context, poolId uint64, statType types.STAKER_STATS, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(util.GetByteKey(string(statType), poolId), bz)
}
