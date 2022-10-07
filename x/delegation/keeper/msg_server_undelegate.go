package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Undelegate handles the transaction of undelegating a given amount from the delegated tokens
// The Undelegation is not performed immediately, instead an unbonding entry is created and pushed
// to a queue. When the unbonding timeout is reached the actual undelegation is performed.
// If the delegator got slashed during the unbonding only the remaining tokens will be returned.
func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Don't allow to undelegate more than currently delegated
	if msg.Amount > k.GetDelegationAmountOfDelegator(ctx, msg.Staker, msg.Creator) {
		return nil, errors.Wrapf(types.ErrNotEnoughDelegation, "")
	}

	// Create and insert Unbonding queue entry.
	k.StartUnbondingDelegator(ctx, msg.Staker, msg.Creator, msg.Amount)

	return &types.MsgUndelegateResponse{}, nil
}
