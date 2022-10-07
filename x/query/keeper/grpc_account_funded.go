package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AccountFundedList(goCtx context.Context, req *types.QueryAccountFundedListRequest) (*types.QueryAccountFundedListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var funded []types.Funded

	for _, pool := range k.poolKeeper.GetAllPools(ctx) {
		funder, _ := pool.GetFunder(req.Address)
		funded = append(funded, types.Funded{
			Amount: funder.Amount,
			Pool: &types.BasicPool{
				Id:         pool.Id,
				Name:       pool.Name,
				Runtime:    pool.Runtime,
				Logo:       pool.Logo,
				TotalFunds: pool.TotalFunds,
				Status:     k.GetPoolStatus(ctx, &pool),
			},
		})

	}

	return &types.QueryAccountFundedListResponse{
		Funded: funded,
	}, nil
}
