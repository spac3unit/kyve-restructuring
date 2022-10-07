package pool

import (
	"github.com/KYVENetwork/chain/x/pool/keeper"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {

	// Set Params
	k.SetParams(ctx, genState.Params)

	// Set all the pool
	for _, elem := range genState.PoolList {
		k.SetPool(ctx, elem)
	}

	// Set pool count
	k.SetPoolCount(ctx, genState.PoolCount)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// Import Params
	genesis.Params = k.GetParams(ctx)

	// Import All Pools
	genesis.PoolList = k.GetAllPools(ctx)

	// Import the number of pools
	genesis.PoolCount = k.GetPoolCount(ctx)

	return genesis
}
