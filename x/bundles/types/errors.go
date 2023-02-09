package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bundles module sentinel errors
var (
	ErrUploaderAlreadyClaimed  = sdkerrors.Register(ModuleName, 1100, "uploader role already claimed")
	ErrInvalidArgs             = sdkerrors.Register(ModuleName, 1107, "invalid args")
	ErrFromIndex               = sdkerrors.Register(ModuleName, 1118, "invalid from index")
	ErrNotDesignatedUploader   = sdkerrors.Register(ModuleName, 1113, "not designated uploader")
	ErrUploadInterval          = sdkerrors.Register(ModuleName, 1108, "upload interval not surpassed")
	ErrMaxBundleSize           = sdkerrors.Register(ModuleName, 1109, "max bundle size was surpassed")
	ErrQuorumNotReached        = sdkerrors.Register(ModuleName, 1111, "no quorum reached")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 1119, "invalid vote %v")
	ErrInvalidStorageId        = sdkerrors.Register(ModuleName, 1120, "current storageId %v does not match provided storageId")
	ErrPoolDisabled            = sdkerrors.Register(ModuleName, 1121, "pool is disabled")
	ErrPoolCurrentlyUpgrading  = sdkerrors.Register(ModuleName, 1122, "pool currently upgrading")
	ErrMinDelegationNotReached = sdkerrors.Register(ModuleName, 1200, "min delegation not reached")
	ErrPoolOutOfFunds          = sdkerrors.Register(ModuleName, 1201, "pool is out of funds")
	ErrBundleDropped           = sdkerrors.Register(ModuleName, 1202, "bundle proposal is dropped")
	ErrAlreadyVotedValid       = sdkerrors.Register(ModuleName, 1204, "already voted valid on bundle proposal")
	ErrAlreadyVotedInvalid     = sdkerrors.Register(ModuleName, 1205, "already voted invalid on bundle proposal")
	ErrAlreadyVotedAbstain     = sdkerrors.Register(ModuleName, 1206, "already voted abstain on bundle proposal")
)
