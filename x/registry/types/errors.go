package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// pool errors
var (
	ErrPoolNotFound = sdkerrors.Register(ModuleName, 1100, "pool with id %v does not exist")
)

// general errors
var (
	ErrFromHeight = sdkerrors.Register(ModuleName, 1118, "invalid from height")
)
