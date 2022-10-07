package keeper

import (
	"encoding/binary"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DoesValaccountExist only checks if the key is present in the KV-Store
// without loading and unmarshalling to full entry
func (k Keeper) DoesValaccountExist(
	ctx sdk.Context,
	poolId uint64,
	stakerAddress string,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefix)
	return store.Has(types.ValaccountKey(poolId, stakerAddress))
}

// GetAllValaccountsOfPool returns a list of all valaccount
func (k Keeper) GetAllValaccountsOfPool(
	ctx sdk.Context,
	poolId uint64,
) (val []*types.Valaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefix)

	iterator := sdk.KVStorePrefixIterator(store, util.GetByteKey(poolId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		valaccount := types.Valaccount{}
		k.cdc.MustUnmarshal(iterator.Value(), &valaccount)
		val = append(val, &valaccount)
	}

	return
}

// GetValaccountsFromStaker returns all pools the staker has valaccounts in
func (k Keeper) GetValaccountsFromStaker(ctx sdk.Context, stakerAddress string) (val []*types.Valaccount) {
	storeIndex2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefixIndex2)

	iterator := sdk.KVStorePrefixIterator(storeIndex2, util.GetByteKey(stakerAddress))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolId := binary.BigEndian.Uint64(iterator.Key()[43 : 43+8])
		valaccount, valaccountFound := k.GetValaccount(ctx, poolId, stakerAddress)

		if valaccountFound {
			val = append(val, &valaccount)
		}
	}

	return val
}

func (k Keeper) GetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) uint64 {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)

	if found {
		return valaccount.Points
	}

	return 0
}

func (k Keeper) AddPoint(ctx sdk.Context, poolId uint64, stakerAddress string) uint64 {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)

	if found {
		valaccount.Points = valaccount.Points + 1
		k.SetValaccount(ctx, valaccount)

		return valaccount.Points
	}

	return 0
}

func (k Keeper) ResetPoints(ctx sdk.Context, poolId uint64, stakerAddress string) {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)

	if found {
		valaccount.Points = 0
		k.SetValaccount(ctx, valaccount)
	}
}

// #############################
// #  Raw KV-Store operations  #
// #############################

// SetValaccount set a specific Valaccount in the store from its index
func (k Keeper) SetValaccount(ctx sdk.Context, valaccount types.Valaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefix)
	b := k.cdc.MustMarshal(&valaccount)
	store.Set(types.ValaccountKey(
		valaccount.PoolId,
		valaccount.Staker,
	), b)

	storeIndex2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefixIndex2)
	storeIndex2.Set(types.ValaccountKeyIndex2(
		valaccount.Staker,
		valaccount.PoolId,
	), []byte{})
}

// removeValaccount removes a Valaccount from the store
func (k Keeper) removeValaccount(ctx sdk.Context, valaccount types.Valaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefix)
	store.Delete(types.ValaccountKey(
		valaccount.PoolId,
		valaccount.Staker,
	))

	storeIndex2 := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefixIndex2)
	storeIndex2.Delete(types.ValaccountKeyIndex2(
		valaccount.Staker,
		valaccount.PoolId,
	))
}

// GetValaccount returns a Valaccount from its index
func (k Keeper) GetValaccount(
	ctx sdk.Context,
	poolId uint64,
	stakerAddress string,
) (val types.Valaccount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefix)

	b := store.Get(types.ValaccountKey(
		poolId,
		stakerAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllValaccounts ...
func (k Keeper) GetAllValaccounts(ctx sdk.Context) (list []types.Valaccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ValaccountPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Valaccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
