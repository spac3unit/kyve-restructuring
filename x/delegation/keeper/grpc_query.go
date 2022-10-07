package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) Slashes(goCtx context.Context, req *types.QuerySlashesRequest) (*types.QuerySlashesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	response := types.QuerySlashesResponse{}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DelegationSlashEntriesKeyPrefix)
	iterator := store.ReverseIterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DelegationSlash
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		response.Slashes = append(response.Slashes, &val)
	}

	return &response, nil
}
