package keeper

import (
	"github.com/KYVENetwork/chain/x/query/types"
)

var _ types.QueryAccountServer = Keeper{}
var _ types.QueryPoolServer = Keeper{}
var _ types.QueryStakersServer = Keeper{}
var _ types.QueryDelegationServer = Keeper{}
var _ types.QueryBundlesServer = Keeper{}
var _ types.QueryParamsServer = Keeper{}
