package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bundles module sentinel errors
var (
	ErrUploaderAlreadyClaimed = sdkerrors.Register(ModuleName, 1100, "uploader role already claimed")
	ErrInvalidArgs            = sdkerrors.Register(ModuleName, 1107, "invalid args")
	ErrFromHeight             = sdkerrors.Register(ModuleName, 1118, "invalid from height")
	ErrToHeight               = sdkerrors.Register(ModuleName, 1123, "invalid to height")
	ErrNotDesignatedUploader  = sdkerrors.Register(ModuleName, 1113, "not designated uploader")
	ErrUploadInterval         = sdkerrors.Register(ModuleName, 1108, "upload interval not surpassed")
	ErrMaxBundleSize          = sdkerrors.Register(ModuleName, 1109, "max bundle size was surpassed")
	ErrAlreadyVoted           = sdkerrors.Register(ModuleName, 1110, "already voted on proposal %v")
	ErrQuorumNotReached       = sdkerrors.Register(ModuleName, 1111, "no quorum reached")
	ErrFromKey                = sdkerrors.Register(ModuleName, 1124, "invalid from key")
	ErrVoterIsUploader        = sdkerrors.Register(ModuleName, 1112, "voter is uploader")
	ErrInvalidVote            = sdkerrors.Register(ModuleName, 1119, "invalid vote %v")
	ErrInvalidStorageId       = sdkerrors.Register(ModuleName, 1120, "current storageId %v does not match provided storageId")
	ErrPoolPaused             = sdkerrors.Register(ModuleName, 1121, "pool is paused")
	ErrPoolCurrentlyUpgrading = sdkerrors.Register(ModuleName, 1122, "pool currently upgrading")
	ErrMinStakeNotReached     = sdkerrors.Register(ModuleName, 1200, "min stake not reached")
	ErrPoolOutOfFunds         = sdkerrors.Register(ModuleName, 1201, "pool is out of funds")
	ErrBundleDropped          = sdkerrors.Register(ModuleName, 1202, "bundle proposal is dropped")
	ErrNoDataBundle           = sdkerrors.Register(ModuleName, 1203, "bundle proposal is of type no data")
	ErrAlreadyVotedValid      = sdkerrors.Register(ModuleName, 1204, "already voted valid on bundle proposal")
	ErrAlreadyVotedInvalid    = sdkerrors.Register(ModuleName, 1205, "already voted invalid on bundle proposal")
	ErrAlreadyVotedAbstain    = sdkerrors.Register(ModuleName, 1206, "already voted abstain on bundle proposal")
)
