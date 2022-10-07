package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// FundPool handles the logic to fund a pool.
// A funder is added to the funders list with the specified amount
// If the funders list is full, it checks if the funder wants to fund
// more than the current lowest funder. If so, the current lowest funder
// will get their tokens back and removed form the funders list.
func (k msgServer) FundPool(goCtx context.Context, msg *types.MsgFundPool) (*types.MsgFundPoolResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, poolFound := k.GetPool(ctx, msg.Id)

	if !poolFound {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Check if funder already exists
	funder, found := pool.GetFunder(msg.Creator)

	if found {
		pool.AddToFunder(funder.Address, msg.Amount)
	} else {
		// If funder does not exist, check if limit is already exceeded.
		if len(pool.Funders) >= types.MaxFunders {
			// If so, check if funder wants to fund more than current lowest funder.
			lowestFunder := pool.Funders[0]
			if msg.Amount > lowestFunder.Amount {
				// Unstake lowest Funder

				err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, lowestFunder.Address, lowestFunder.Amount)
				if err != nil {
					return nil, err
				}

				// Emit a defund event.
				if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventDefundPool{
					PoolId:  msg.Id,
					Address: lowestFunder.Address,
					Amount:  lowestFunder.Amount,
				}); errEmit != nil {
					return nil, errEmit
				}

				// Remove from pool
				pool.RemoveFunder(*pool.Funders[0])
			} else {
				return nil, sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrFundsTooLow.Error(), lowestFunder.Amount)
			}
		}

		pool.InsertFunder(types.Funder{
			Address: msg.Creator,
			Amount:  msg.Amount,
		})
	}

	err := util.TransferFromAddressToModule(k.bankKeeper, ctx, msg.Creator, types.ModuleName, msg.Amount)
	if err != nil {
		return nil, err
	}

	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventFundPool{
		PoolId:  msg.Id,
		Address: msg.Creator,
		Amount:  msg.Amount,
	}); errEmit != nil {
		return nil, errEmit
	}

	k.SetPool(ctx, pool)

	return &types.MsgFundPoolResponse{}, nil
}
