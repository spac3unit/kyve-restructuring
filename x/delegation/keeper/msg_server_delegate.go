package keeper

import (
	"context"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KYVENetwork/chain/util"

	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Delegate handles the transaction of delegating a specific amount of $KYVE to a staker
// The only requirement for the transaction to succeed is that the staker exists
// and the user has enough balance.
func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.stakersKeeper.DoesStakerExist(ctx, msg.Staker) {
		return nil, sdkErrors.Wrap(types.ErrStakerDoesNotExist, msg.Staker)
	}

	// Performs logical delegation without transferring the amount
	k.performDelegation(ctx, msg.Staker, msg.Creator, msg.Amount)

	// Transfer tokens from sender to this module.
	if transferErr := util.TransferFromAddressToModule(k.bankKeeper, ctx, msg.Creator, types.ModuleName, msg.Amount); transferErr != nil {
		return nil, transferErr
	}

	// Emit a delegation event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventDelegate{
		Address: msg.Creator,
		Node:    msg.Staker,
		Amount:  msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgDelegateResponse{}, nil
}
