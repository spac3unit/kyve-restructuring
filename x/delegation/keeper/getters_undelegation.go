package keeper

import (
	"encoding/binary"
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// #####################
// === QUEUE ENTRIES ===
// #####################

// SetUndelegationQueueEntry ...
func (k Keeper) SetUndelegationQueueEntry(ctx sdk.Context, undelegationQueueEntry types.UndelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UndelegationQueueKeyPrefix)
	b := k.cdc.MustMarshal(&undelegationQueueEntry)
	store.Set(types.UndelegationQueueKey(
		undelegationQueueEntry.Index,
	), b)

	// Insert the same entry with a different key prefix for query lookup
	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.UndelegationQueueKeyPrefixIndex2)
	indexStore.Set(types.UndelegationQueueKeyIndex2(
		undelegationQueueEntry.Delegator,
		undelegationQueueEntry.Index,
	), []byte{})
}

// GetUndelegationQueueEntry ...
func (k Keeper) GetUndelegationQueueEntry(ctx sdk.Context, index uint64) (val types.UndelegationQueueEntry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UndelegationQueueKeyPrefix)

	b := store.Get(types.UndelegationQueueKey(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveUndelegationQueueEntry ...
func (k Keeper) RemoveUndelegationQueueEntry(ctx sdk.Context, undelegationQueueEntry *types.UndelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UndelegationQueueKeyPrefix)
	store.Delete(types.UndelegationQueueKey(undelegationQueueEntry.Index))

	indexStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.UndelegationQueueKeyPrefixIndex2)
	indexStore.Delete(types.UndelegationQueueKeyIndex2(
		undelegationQueueEntry.Delegator,
		undelegationQueueEntry.Index,
	))
}

// GetAllUnbondingDelegationQueueEntries returns all delegator unbondings
func (k Keeper) GetAllUnbondingDelegationQueueEntries(ctx sdk.Context) (list []types.UndelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.UndelegationQueueKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.UndelegationQueueEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllUnbondingDelegationQueueEntriesOfDelegator returns all delegator unbondings of the given address
func (k Keeper) GetAllUnbondingDelegationQueueEntriesOfDelegator(ctx sdk.Context, address string) (list []types.UndelegationQueueEntry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), util.GetByteKey(types.UndelegationQueueKeyPrefixIndex2, address))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		index := binary.BigEndian.Uint64(iterator.Key()[0:8])

		entry, _ := k.GetUndelegationQueueEntry(ctx, index)
		list = append(list, entry)
	}

	return
}

// ###################
// === QUEUE STATE ===
// ###################

// GetQueueState returns the state for the undelegation queue
func (k Keeper) GetQueueState(ctx sdk.Context) (state types.QueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := store.Get(types.QueueKey)

	if b == nil {
		return state
	}

	k.cdc.MustUnmarshal(b, &state)
	return
}

// SetQueueState saves the undelegation queue state
func (k Keeper) SetQueueState(ctx sdk.Context, state types.QueueState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	b := k.cdc.MustMarshal(&state)
	store.Set(types.QueueKey, b)
}
