package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Stakers(c context.Context, req *types.QueryStakersRequest) (*types.QueryStakersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	data := make([]types.FullStaker, 0)
	pageRes, err := k.delegationKeeper.GetPaginatedStakersByDelegation(ctx, req.Pagination, func(address string) {
		data = append(data, *k.GetFullStaker(ctx, address))
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryStakersResponse{Stakers: data, Pagination: pageRes}, nil
}

func (k Keeper) Staker(c context.Context, req *types.QueryStakerRequest) (*types.QueryStakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if !k.stakerKeeper.DoesStakerExist(ctx, req.Address) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryStakerResponse{Staker: *k.GetFullStaker(ctx, req.Address)}, nil
}
