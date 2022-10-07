package types

import (
	"github.com/KYVENetwork/chain/util"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyVoteSlash            = []byte("VoteSlash")
	DefaultVoteSlash string = "0.1"
)

var (
	KeyUploadSlash            = []byte("UploadSlash")
	DefaultUploadSlash string = "0.2"
)

var (
	KeyTimeoutSlash            = []byte("TimeoutSlash")
	DefaultTimeoutSlash string = "0.02"
)

var (
	KeyUnbondingStakingTime            = []byte("UnbondingStakingTime")
	DefaultUnbondingStakingTime uint64 = 60 * 60 * 24 * 5
)

var (
	KeyCommissionChangeTime            = []byte("CommissionChangeTime")
	DefaultCommissionChangeTime uint64 = 60 * 60 * 24 * 5
)

var (
	KeyLeavePoolTime            = []byte("LeavePoolTime")
	DefaultLeavePoolTime uint64 = 60 * 60 * 24 * 5
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	voteSlash string,
	uploadSlash string,
	timeoutSlash string,
	unbondingStakingTime uint64,
	commissionChangeTime uint64,
	leavePoolTime uint64,
) Params {
	return Params{
		VoteSlash:            voteSlash,
		UploadSlash:          uploadSlash,
		TimeoutSlash:         timeoutSlash,
		UnbondingStakingTime: unbondingStakingTime,
		CommissionChangeTime: commissionChangeTime,
		LeavePoolTime:        leavePoolTime,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultVoteSlash,
		DefaultUploadSlash,
		DefaultTimeoutSlash,
		DefaultUnbondingStakingTime,
		DefaultCommissionChangeTime,
		DefaultLeavePoolTime,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyVoteSlash, &p.VoteSlash, util.ValidatePercentage),
		paramtypes.NewParamSetPair(KeyUploadSlash, &p.UploadSlash, util.ValidatePercentage),
		paramtypes.NewParamSetPair(KeyTimeoutSlash, &p.TimeoutSlash, util.ValidatePercentage),
		paramtypes.NewParamSetPair(KeyUnbondingStakingTime, &p.UnbondingStakingTime, util.ValidateUint64),
		paramtypes.NewParamSetPair(KeyCommissionChangeTime, &p.CommissionChangeTime, util.ValidateUint64),
		paramtypes.NewParamSetPair(KeyLeavePoolTime, &p.LeavePoolTime, util.ValidateUint64),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := util.ValidatePercentage(p.VoteSlash); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.UploadSlash); err != nil {
		return err
	}

	if err := util.ValidatePercentage(p.TimeoutSlash); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.UnbondingStakingTime); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.CommissionChangeTime); err != nil {
		return err
	}

	if err := util.ValidateUint64(p.LeavePoolTime); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
