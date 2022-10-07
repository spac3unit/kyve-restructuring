package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/KYVENetwork/chain/x/delegation/types"
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

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		distrKeeper   types.DistrKeeper
		poolKeeper    types.PoolKeeper
		upgradeKeeper types.UpgradeKeeper
		stakersKeeper types.StakersKeeper
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
	upgradeKeeper types.UpgradeKeeper,
	stakersKeeper types.StakersKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distrkeeper,
		poolKeeper:    poolKeeper,
		upgradeKeeper: upgradeKeeper,
		stakersKeeper: stakersKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) StoreKey() sdk.StoreKey {
	return k.storeKey
}

var memStoreInitialized = false

func (k Keeper) InitMemStore(ctx sdk.Context) {
	if !memStoreInitialized {

		// Update mem index
		for _, entry := range k.GetAllDelegationData(ctx) {
			k.SetStakerIndex(ctx, entry.Staker)
		}

		memStoreInitialized = true
	}
}
