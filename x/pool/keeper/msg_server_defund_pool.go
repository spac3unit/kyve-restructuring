package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// DefundPool handles the logic to defund a pool.
// If the user is a funder, it will subtract the provided amount
// and send the tokens back. If the amount equals the current funding amount
// the funder is removed completely.
func (k msgServer) DefundPool(goCtx context.Context, msg *types.MsgDefundPool) (*types.MsgDefundPoolResponse, error) {

	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if the sender is a funder in this pool.
	funder, isFunder := pool.GetFunder(msg.Creator)
	if !isFunder {
		return nil, sdkErrors.ErrNotFound
	}

	// Check if the sender is trying to defund more than they have funded.
	if msg.Amount > funder.Amount {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrDefundTooHigh.Error(), funder.Address)
	}

	// Update state variables (or completely remove if fully defunding).
	pool.SubFromFunder(funder.Address, msg.Amount)

	// Transfer tokens from this module to sender.
	if err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Creator, msg.Amount); err != nil {
		return nil, err
	}

	// Emit a defund event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
		PoolId:  msg.Id,
		Address: msg.Creator,
		Amount:  msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	k.SetPool(ctx, pool)

	return &types.MsgDefundPoolResponse{}, nil
}
