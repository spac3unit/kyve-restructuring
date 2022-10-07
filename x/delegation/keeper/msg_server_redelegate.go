package keeper

import (
	"context"
	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Redelegate lets a user redelegate from one staker to another staker
// The user has N redelegation spells. When this transaction is executed
// one spell is used. When all spells are consumed the transaction fails.
// The user then needs to wait for the oldest spell to expire to call
// this transaction again.
func (k msgServer) Redelegate(goCtx context.Context, msg *types.MsgRedelegate) (*types.MsgRedelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender is a delegator
	if !k.DoesDelegatorExist(ctx, msg.FromStaker, msg.Creator) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNotADelegator.Error())
	}

	// Check if destination staker exists
	if !k.stakersKeeper.DoesStakerExist(ctx, msg.ToStaker) {
		return nil, types.ErrStakerDoesNotExist
	}

	// Check if the sender is trying to undelegate more than he has delegated.
	if msg.Amount > k.GetDelegationAmountOfDelegator(ctx, msg.FromStaker, msg.Creator) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrInsufficientFunds, types.ErrNotEnoughDelegation.Error(), msg.Amount)
	}

	// Only errors if all spells are currently on cooldown
	if err := k.consumeRedelegationSpell(ctx, msg.Creator); err != nil {
		return nil, err
	}

	// The redelegation is translated into an undelegation from the old staker ...
	if actualAmount := k.performUndelegation(ctx, msg.FromStaker, msg.Creator, msg.Amount); actualAmount != msg.Amount {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrInsufficientFunds, types.ErrNotEnoughDelegation.Error(), msg.Amount)
	}
	// ... and a new delegation to the new staker
	k.performDelegation(ctx, msg.ToStaker, msg.Creator, msg.Amount)

	// Emit a delegation event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventRedelegate{
		Address:  msg.Creator,
		FromNode: msg.FromStaker,
		ToNode:   msg.ToStaker,
		Amount:   msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgRedelegateResponse{}, nil
}
