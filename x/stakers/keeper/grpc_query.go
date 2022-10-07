package keeper

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
)

var _ types.QueryServer = Keeper{}
