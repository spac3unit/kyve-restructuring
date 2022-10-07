package types

import (
	"github.com/KYVENetwork/chain/util"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyUnbondingDelegationTime            = []byte("UnbondingDelegationTime")
	DefaultUnbondingDelegationTime uint64 = 60 * 60 * 24 * 5
)

var (
	KeyRedelegationCooldown            = []byte("RedelegationCooldown")
	DefaultRedelegationCooldown uint64 = 60 * 60 * 24 * 5
)

var (
	KeyRedelegationMaxAmount            = []byte("RedelegationMaxAmount")
	DefaultRedelegationMaxAmount uint64 = 5
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	unbondingDelegationTime uint64,
	redelegationCooldown uint64,
	redelegationMaxAmount uint64,
) Params {
	return Params{
		UnbondingDelegationTime: unbondingDelegationTime,
		RedelegationCooldown:    redelegationCooldown,
		RedelegationMaxAmount:   redelegationMaxAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultUnbondingDelegationTime,
		DefaultRedelegationCooldown,
		DefaultRedelegationMaxAmount,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyUnbondingDelegationTime, &p.UnbondingDelegationTime, util.ValidateUint64),
		paramtypes.NewParamSetPair(KeyRedelegationCooldown, &p.RedelegationCooldown, util.ValidateUint64),
		paramtypes.NewParamSetPair(KeyRedelegationMaxAmount, &p.RedelegationMaxAmount, util.ValidateUint64),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidateUint64(p.UnbondingDelegationTime); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.RedelegationCooldown); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.RedelegationMaxAmount); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
