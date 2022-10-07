package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) LeavePool(goCtx context.Context, msg *types.MsgLeavePool) (*types.MsgLeavePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valaccount, valaccountFound := k.GetValaccount(ctx, msg.PoolId, msg.Creator)
	if !valaccountFound {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrInvalidRequest, types.ErrAlreadyLeftPool.Error())
	}

	valaccount.IsLeaving = true
	k.SetValaccount(ctx, valaccount)

	err := k.orderLeavePool(ctx, msg.Creator, msg.PoolId)
	if err != nil {
		return nil, err
	}

	return &types.MsgLeavePoolResponse{}, nil
}
