package keeper

import (
	"github.com/KYVENetwork/chain/util"
	"github.com/KYVENetwork/chain/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetBondingOfAddress(ctx sdk.Context, address sdk.AccAddress) sdk.Int {

	var amount uint64 = 0

	delegatorStore := prefix.NewStore(ctx.KVStore(k.storeKey), util.GetByteKey(types.DelegatorKeyPrefixIndex2, address.String()))
	delegatorIterator := sdk.KVStorePrefixIterator(delegatorStore, nil)
	defer delegatorIterator.Close()

	for ; delegatorIterator.Valid(); delegatorIterator.Next() {
		stakerAddress := string(delegatorIterator.Key()[0:43])
		amount += k.GetDelegationAmountOfDelegator(ctx, stakerAddress, address.String())
	}

	return sdk.NewIntFromUint64(amount)
}

func (k Keeper) TotalProtocolBonding(ctx sdk.Context) sdk.Int {
	total := uint64(0)

	// TODO consider aggregated variable
	for _, data := range k.GetAllDelegationData(ctx) {
		total += data.TotalDelegation
	}

	return sdk.NewIntFromUint64(total)
}
