package types

import (
	"github.com/KYVENetwork/chain/util"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyUploadTimeout            = []byte("UploadTimeout")
	DefaultUploadTimeout uint64 = 600
)

var (
	KeyStorageCost            = []byte("StorageCost")
	DefaultStorageCost uint64 = 100000
)

var (
	KeyNetworkFee            = []byte("NetworkFee")
	DefaultNetworkFee string = "0.01"
)

var (
	KeyMaxPoints            = []byte("MaxPoints")
	DefaultMaxPoints uint64 = 5
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	uploadTimeout uint64,
	storageCost uint64,
	networkFee string,
	maxPoints uint64,
) Params {
	return Params{
		UploadTimeout: uploadTimeout,
		StorageCost:   storageCost,
		NetworkFee:    networkFee,
		MaxPoints:     maxPoints,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultUploadTimeout,
		DefaultStorageCost,
		DefaultNetworkFee,
		DefaultMaxPoints,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyUploadTimeout, &p.UploadTimeout, util.ValidateUint64),
		paramtypes.NewParamSetPair(KeyStorageCost, &p.StorageCost, util.ValidateUint64),
		paramtypes.NewParamSetPair(KeyNetworkFee, &p.NetworkFee, util.ValidatePercentage),
		paramtypes.NewParamSetPair(KeyMaxPoints, &p.MaxPoints, util.ValidateUint64),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {

	if err := util.ValidateUint64(p.UploadTimeout); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.StorageCost); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.NetworkFee); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.MaxPoints); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
