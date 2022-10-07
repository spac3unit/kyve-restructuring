package keeper

import (
	"encoding/json"
	"github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CreatePool(ctx sdk.Context, p *types.CreatePoolProposal) error {
	// Validate config json
	if !json.Valid([]byte(p.Config)) {
		return sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidJson.Error(), p.Config)
	}

	// Validate binaries json
	if !json.Valid([]byte(p.Binaries)) {
		return sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidJson.Error(), p.Binaries)
	}

	k.AppendPool(ctx, types.Pool{
		Name:           p.Name,
		Runtime:        p.Runtime,
		Logo:           p.Logo,
		Config:         p.Config,
		StartKey:       p.StartKey,
		UploadInterval: p.UploadInterval,
		OperatingCost:  p.OperatingCost,
		MinStake:       p.MinStake,
		MaxBundleSize:  p.MaxBundleSize,
		Protocol: &types.Protocol{
			Version:     p.Version,
			Binaries:    p.Binaries,
			LastUpgrade: uint64(ctx.BlockTime().Unix()),
		},
		UpgradePlan: &types.UpgradePlan{},
	})

	if errEmit := ctx.EventManager().EmitTypedEvent(&types.EventCreatePool{
		Id:             k.GetPoolCount(ctx),
		Name:           p.Name,
		Runtime:        p.Runtime,
		Logo:           p.Logo,
		Config:         p.Config,
		StartKey:       p.StartKey,
		UploadInterval: p.UploadInterval,
		OperatingCost:  p.OperatingCost,
		MinStake:       p.MinStake,
		MaxBundleSize:  p.MaxBundleSize,
		Version:        p.Version,
		Binaries:       p.Binaries,
	}); errEmit != nil {
		return errEmit
	}

	return nil
}

func (k Keeper) UpdatePool(ctx sdk.Context, p *types.UpdatePoolProposal) error {
	pool, found := k.GetPool(ctx, p.Id)
	if !found {
		return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), p.Id)
	}

	type Update struct {
		Name           *string
		Runtime        *string
		Logo           *string
		Config         *string
		UploadInterval *uint64
		OperatingCost  *uint64
		MinStake       *uint64
		MaxBundleSize  *uint64
	}

	var update Update

	if err := json.Unmarshal([]byte(p.Payload), &update); err != nil {
		return err
	}

	if update.Name != nil {
		pool.Name = *update.Name
	}

	if update.Runtime != nil {
		pool.Runtime = *update.Runtime
	}

	if update.Logo != nil {
		pool.Logo = *update.Logo
	}

	if update.Config != nil {
		if json.Valid([]byte(*update.Config)) {
			pool.Config = *update.Config
		} else {
			return sdkErrors.Wrapf(sdkErrors.ErrLogic, types.ErrInvalidJson.Error(), *update.Config)
		}
	}

	if update.UploadInterval != nil {
		pool.UploadInterval = *update.UploadInterval
	}

	if update.OperatingCost != nil {
		pool.OperatingCost = *update.OperatingCost
	}

	if update.MinStake != nil {
		pool.MinStake = *update.MinStake
	}

	if update.MaxBundleSize != nil {
		pool.MaxBundleSize = *update.MaxBundleSize
	}

	k.SetPool(ctx, pool)

	return nil
}

func (k Keeper) PausePool(ctx sdk.Context, p *types.PausePoolProposal) error {
	// Attempt to fetch the pool, throw an error if not found.
	pool, found := k.GetPool(ctx, p.Id)
	if !found {
		return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), p.Id)
	}

	// Throw an error if the pool is already paused.
	if pool.Paused {
		return sdkErrors.Wrapf(sdkErrors.ErrLogic, "Pool is already paused.")
	}

	// Pause the pool and return.
	pool.Paused = true
	k.SetPool(ctx, pool)

	return nil
}

func (k Keeper) UnpausePool(ctx sdk.Context, p *types.UnpausePoolProposal) error {
	// Attempt to fetch the pool, throw an error if not found.
	pool, found := k.GetPool(ctx, p.Id)
	if !found {
		return sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), p.Id)
	}

	// Throw an error if the pool is already unpaused.
	if !pool.Paused {
		return sdkErrors.Wrapf(sdkErrors.ErrLogic, "Pool is already unpaused.")
	}

	// Unpause the pool and return.
	pool.Paused = false
	k.SetPool(ctx, pool)

	return nil
}

func (k Keeper) UpgradePool(ctx sdk.Context, p *types.SchedulePoolUpgradeProposal) error {
	// Check if upgrade version and binaries are not empty
	if p.Version == "" || p.Binaries == "" {
		return types.ErrInvalidArgs
	}

	var scheduledAt uint64

	// If upgrade time was already surpassed we upgrade immediately
	if p.ScheduledAt < uint64(ctx.BlockTime().Unix()) {
		scheduledAt = uint64(ctx.BlockTime().Unix())
	} else {
		scheduledAt = p.ScheduledAt
	}

	// go through every pool and schedule the upgrade
	for _, pool := range k.GetAllPools(ctx) {
		// Skip if runtime does not match
		if pool.Runtime != p.Runtime {
			continue
		}

		// Skip if pool is currently upgrading
		if pool.UpgradePlan.ScheduledAt > 0 {
			continue
		}

		// register upgrade plan
		pool.UpgradePlan = &types.UpgradePlan{
			Version:     p.Version,
			Binaries:    p.Binaries,
			ScheduledAt: scheduledAt,
			Duration:    p.Duration,
		}

		// Update the pool
		k.SetPool(ctx, pool)
	}

	return nil
}

func (k Keeper) CancelPoolUpgrade(ctx sdk.Context, p *types.CancelPoolUpgradeProposal) error {
	// go through every pool and cancel the upgrade
	for _, pool := range k.GetAllPools(ctx) {
		// Skip if runtime does not match
		if pool.Runtime != p.Runtime {
			continue
		}

		// Continue if there is no upgrade scheduled
		if pool.UpgradePlan.ScheduledAt == 0 {
			continue
		}

		// clear upgrade plan
		pool.UpgradePlan = &types.UpgradePlan{}

		// Update the pool
		k.SetPool(ctx, pool)
	}

	return nil
}
