package pool

import (
	"github.com/KYVENetwork/chain/x/pool/keeper"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func NewPoolProposalHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CreatePoolProposal:
			return k.CreatePool(ctx, c)
		case *types.UpdatePoolProposal:
			return k.UpdatePool(ctx, c)
		case *types.PausePoolProposal:
			return k.PausePool(ctx, c)
		case *types.UnpausePoolProposal:
			return k.UnpausePool(ctx, c)
		case *types.SchedulePoolUpgradeProposal:
			return k.UpgradePool(ctx, c)
		case *types.CancelPoolUpgradeProposal:
			return k.CancelPoolUpgrade(ctx, c)

		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized proposal content type: %T", c)
		}
	}
}
