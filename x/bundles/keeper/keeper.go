package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper    types.AccountKeeper
		bankKeeper       types.BankKeeper
		distrkeeper      types.DistrKeeper
		poolKeeper       types.PoolKeeper
		stakerKeeper     types.StakerKeeper
		delegationKeeper types.DelegationKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	distrkeeper types.DistrKeeper,
	poolKeeper types.PoolKeeper,
	stakerKeeper types.StakerKeeper,
	delegationKeeper types.DelegationKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		accountKeeper:    accountKeeper,
		bankKeeper:       bankKeeper,
		distrkeeper:      distrkeeper,
		poolKeeper:       poolKeeper,
		stakerKeeper:     stakerKeeper,
		delegationKeeper: delegationKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

var memStoreInitialized = false

func (k Keeper) InitMemStore(ctx sdk.Context) {
	if !memStoreInitialized {

		// Update mem index
		for _, entry := range k.GetAllFinalizedBundles(ctx) {
			k.SetFinalizedBundleIndexes(ctx, entry)
		}

		memStoreInitialized = true
	}
}
