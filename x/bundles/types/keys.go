package types

import (
	"encoding/binary"
	"github.com/KYVENetwork/chain/util"
)

const (
	// ModuleName defines the module name
	ModuleName = "bundles"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bundles"
)

var (
	// BundleKeyPrefix ...
	BundleKeyPrefix = []byte{1}
	// FinalizedBundlePrefix ...
	FinalizedBundlePrefix = []byte{2}

	FinalizedBundleByStorageIdPrefix = []byte{10}
	FinalizedBundleByHeightPrefix    = []byte{11}
)

// BundleProposalKey ...
func BundleProposalKey(poolId uint64) []byte {
	return util.GetByteKey(poolId)
}

// FinalizedBundleKey ...
func FinalizedBundleKey(poolId uint64, id uint64) []byte {
	return util.GetByteKey(poolId, id)
}

// FinalizedBundleByStorageIdKey ...
func FinalizedBundleByStorageIdKey(storageId string) []byte {
	return util.GetByteKey(storageId)
}

func FinalizedBundlePoolIdAndIdToByte(poolId uint64, id uint64) []byte {
	return util.GetByteKey(poolId, id)
}

func FinalizedBundlePoolIdAndIdToValue(bytes []byte) (poolId uint64, id uint64) {
	poolId = binary.BigEndian.Uint64(bytes[0:8])
	id = binary.BigEndian.Uint64(bytes[8:16])
	return
}

// FinalizedBundleByHeightKey ...
func FinalizedBundleByHeightKey(poolId uint64, height uint64) []byte {
	return util.GetByteKey(poolId, height)
}
