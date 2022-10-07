package keeper

import (
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
)

func (k Keeper) SetStakerIndex(ctx sdk.Context, staker string) {
	amount := k.GetDelegationAmount(ctx, staker)
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)
	store.Set(types.StakerIndexKey(math.MaxUint64-amount, staker), []byte{0})
}

func (k Keeper) RemoveStakerIndex(ctx sdk.Context, staker string) {
	amount := k.GetDelegationAmount(ctx, staker)
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)
	store.Delete(types.StakerIndexKey(math.MaxUint64-amount, staker))
}

func (k Keeper) GetPaginatedStakersByDelegation(ctx sdk.Context, pagination *query.PageRequest, accumulator func(staker string)) (*query.PageResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.memKey), types.StakerIndexKeyPrefix)

	pageRes, err := query.FilteredPaginate(store, pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			address := string(key[8 : 8+43])
			accumulator(address)
		}
		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return pageRes, nil
}
