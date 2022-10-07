package keeper

import (
	"fmt"

	bundlekeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	delegationkeeper "github.com/KYVENetwork/chain/x/delegation/keeper"
	poolkeeper "github.com/KYVENetwork/chain/x/pool/keeper"
	stakerskeeper "github.com/KYVENetwork/chain/x/stakers/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/KYVENetwork/chain/x/query/types"
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

		accountKeeper    authkeeper.AccountKeeper
		bankKeeper       bankkeeper.Keeper
		distrkeeper      distrkeeper.Keeper
		poolKeeper       poolkeeper.Keeper
		stakerKeeper     stakerskeeper.Keeper
		delegationKeeper delegationkeeper.Keeper
		bundleKeeper     bundlekeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	distrkeeper distrkeeper.Keeper,
	poolKeeper poolkeeper.Keeper,
	stakerKeeper stakerskeeper.Keeper,
	delegationKeeper delegationkeeper.Keeper,
	bundleKeeper bundlekeeper.Keeper,
) *Keeper {

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
		bundleKeeper:     bundleKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
