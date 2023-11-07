package types

import (
	"cosmossdk.io/errors"
)

// x/bundles module sentinel errors
var (
	ErrUploaderAlreadyClaimed  = errors.Register(ModuleName, 1100, "uploader role already claimed")
	ErrInvalidArgs             = errors.Register(ModuleName, 1107, "invalid args")
	ErrFromIndex               = errors.Register(ModuleName, 1118, "invalid from index")
	ErrNotDesignatedUploader   = errors.Register(ModuleName, 1113, "not designated uploader")
	ErrUploadInterval          = errors.Register(ModuleName, 1108, "upload interval not surpassed")
	ErrMaxBundleSize           = errors.Register(ModuleName, 1109, "max bundle size was surpassed")
	ErrQuorumNotReached        = errors.Register(ModuleName, 1111, "no quorum reached")
	ErrInvalidVote             = errors.Register(ModuleName, 1119, "invalid vote %v")
	ErrInvalidStorageId        = errors.Register(ModuleName, 1120, "current storageId %v does not match provided storageId")
	ErrPoolDisabled            = errors.Register(ModuleName, 1121, "pool is disabled")
	ErrPoolCurrentlyUpgrading  = errors.Register(ModuleName, 1122, "pool currently upgrading")
	ErrMinDelegationNotReached = errors.Register(ModuleName, 1200, "min delegation not reached")
	ErrBundleDropped           = errors.Register(ModuleName, 1202, "bundle proposal is dropped")
	ErrAlreadyVotedValid       = errors.Register(ModuleName, 1204, "already voted valid on bundle proposal")
	ErrAlreadyVotedInvalid     = errors.Register(ModuleName, 1205, "already voted invalid on bundle proposal")
	ErrAlreadyVotedAbstain     = errors.Register(ModuleName, 1206, "already voted abstain on bundle proposal")
	ErrVotingPowerTooHigh      = errors.Register(ModuleName, 1207, "staker in pool has more than 50% voting power")
)
