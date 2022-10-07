package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/util"

	"github.com/KYVENetwork/chain/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// WithdrawRewards calculates the current rewards of a delegator and transfers to balance to
// the delegator's wallet. Only the delegator himself can call this transaction.
func (k msgServer) WithdrawRewards(goCtx context.Context, msg *types.MsgWithdrawRewards) (*types.MsgWithdrawRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the sender has delegated to the given staker
	if !k.DoesDelegatorExist(ctx, msg.Staker, msg.Creator) {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrNotADelegator.Error())
	}

	// Withdraw all rewards of the sender.
	reward := k.f1WithdrawRewards(ctx, msg.Staker, msg.Creator)

	// Transfer reward $KYVE from this module to sender.
	err := util.TransferFromModuleToAddress(k.bankKeeper, ctx, types.ModuleName, msg.Creator, reward)
	if err != nil {
		return nil, err
	}

	// Emit a delegation event.
	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventWithdrawRewards{
		Address:  msg.Creator,
		FromNode: msg.Staker,
		Amount:   reward,
	}); errEmit != nil {
		return nil, errEmit
	}

	return &types.MsgWithdrawRewardsResponse{}, nil
}
