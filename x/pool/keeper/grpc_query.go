package keeper

import (
	"github.com/KYVENetwork/chain/x/pool/types"
)

var _ types.QueryServer = Keeper{}
