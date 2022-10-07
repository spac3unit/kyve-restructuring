package types

import (
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "stakers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_stakers"
)

var (
	// StakerKeyPrefix is the prefix to retrieve all Staker
	StakerKeyPrefix = []byte{1}

	ValaccountPrefix       = []byte{2, 0}
	ValaccountPrefixIndex2 = []byte{2, 1}

	// UnbondingStakingEntryKeyPrefix ...
	UnbondingStakingEntryKeyPrefix = []byte{3, 0}
	// UnbondingStakingEntryKeyPrefixIndex2 ...
	UnbondingStakingEntryKeyPrefixIndex2 = []byte{3, 1}

	// CommissionChangeEntryKeyPrefix ...
	CommissionChangeEntryKeyPrefix = []byte{4, 0}
	// CommissionChangeEntryKeyPrefixIndex2 ...
	CommissionChangeEntryKeyPrefixIndex2 = []byte{4, 1}

	// LeavePoolEntryKeyPrefix ...
	LeavePoolEntryKeyPrefix = []byte{5, 0}
	// LeavePoolEntryKeyPrefixIndex2 ...
	LeavePoolEntryKeyPrefixIndex2 = []byte{5, 1}
)

// ENUM aggregated data types
type STAKER_STATS string

var (
	STAKER_STATS_TOTAL_STAKE STAKER_STATS = "total_stake"
	STAKER_STATS_COUNT       STAKER_STATS = "total_stakers"
)

// ENUM queue types identifiers
type QUEUE_IDENTIFIER []byte

var (
	QUEUE_IDENTIFIER_COMMISSION QUEUE_IDENTIFIER = []byte{30, 2}
	QUEUE_IDENTIFIER_LEAVE      QUEUE_IDENTIFIER = []byte{30, 3}
)

const (
	MaxStakers        = 50
	DefaultCommission = "0.9"
)

// StakerKey returns the store Key to retrieve a Staker from the index fields
func StakerKey(staker string) []byte {
	return util.GetByteKey(staker)
}

func ValaccountKey(poolId uint64, staker string) []byte {
	return util.GetByteKey(poolId, staker)
}

func ValaccountKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}

func CommissionChangeEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}

// Important: only one queue entry per staker+poolId is allowed at a time.
func CommissionChangeEntryKeyIndex2(staker string) []byte {
	return util.GetByteKey(staker)
}

func UnbondingStakeEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}
func UnbondingStakeEntryKeyIndex2(staker string, index uint64) []byte {
	return util.GetByteKey(staker, index)
}

func LeavePoolEntryKey(index uint64) []byte {
	return util.GetByteKey(index)
}
func LeavePoolEntryKeyIndex2(staker string, poolId uint64) []byte {
	return util.GetByteKey(staker, poolId)
}
