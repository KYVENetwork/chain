package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// pool errors
var (
	ErrPoolNotFound = sdkerrors.Register(ModuleName, 1100, "pool with id %v does not exist")
)

// funding errors
var (
	ErrFundsTooLow   = sdkerrors.Register(ModuleName, 1101, "minimum funding amount of %vkyve not reached")
	ErrDefundTooHigh = sdkerrors.Register(ModuleName, 1102, "maximum defunding amount of %vkyve surpassed")
)

// staking errors
var (
	ErrStakeTooLow    = sdkerrors.Register(ModuleName, 1103, "minimum staking amount of %vkyve not reached")
	ErrUnstakeTooHigh = sdkerrors.Register(ModuleName, 1104, "maximum unstaking amount of %vkyve surpassed")
)

// general errors
var (
	ErrNoStaker               = sdkerrors.Register(ModuleName, 1105, "sender is no staker")
	ErrPoolPaused             = sdkerrors.Register(ModuleName, 1106, "pool is paused")
	ErrInvalidArgs            = sdkerrors.Register(ModuleName, 1107, "invalid args")
	ErrUploadInterval         = sdkerrors.Register(ModuleName, 1108, "upload interval not surpassed")
	ErrInvalidStorageId       = sdkerrors.Register(ModuleName, 1109, "current storageId %v does not match provided storageId")
	ErrAlreadyVoted           = sdkerrors.Register(ModuleName, 1110, "already voted on proposal %v")
	ErrQuorumNotReached       = sdkerrors.Register(ModuleName, 1111, "quorum not reached")
	ErrVoterIsUploader        = sdkerrors.Register(ModuleName, 1112, "voter is uploader")
	ErrNotDesignatedUploader  = sdkerrors.Register(ModuleName, 1113, "not designated uploader")
	ErrUploaderAlreadyClaimed = sdkerrors.Register(ModuleName, 1114, "uploader role already claimed")
	ErrNotEnoughNodesOnline   = sdkerrors.Register(ModuleName, 1115, "not enough nodes online")
	ErrInvalidCommission      = sdkerrors.Register(ModuleName, 1116, "invalid commission %v")
	ErrSelfDelegation         = sdkerrors.Register(ModuleName, 1117, "self delegation not allowed")
	ErrFromHeight             = sdkerrors.Register(ModuleName, 1118, "invalid from height")
	ErrInvalidVote            = sdkerrors.Register(ModuleName, 1119, "invalid vote %v")
	ErrMaxBundleSize          = sdkerrors.Register(ModuleName, 1120, "bundle size is too high")
	ErrPoolCurrentlyUpgrading = sdkerrors.Register(ModuleName, 1121, "pool currently upgrading")
	ErrPoolNoUpgradeScheduled = sdkerrors.Register(ModuleName, 1122, "pool has no scheduled upgrade")
	ErrToHeight               = sdkerrors.Register(ModuleName, 1123, "invalid to height")
	ErrFromKey                = sdkerrors.Register(ModuleName, 1124, "invalid from key")
	ErrProposalNotFound       = sdkerrors.Register(ModuleName, 1125, "proposal with pool id %v and bundle id %v does not exist")
	ErrNotEnoughStake   = sdkerrors.Register(ModuleName, 1126, "not enough stake in pool")
)

// delegation errors
var (
	ErrNotADelegator                   = sdkerrors.Register(ModuleName, 1127, "not a delegator")
	ErrNotEnoughDelegation             = sdkerrors.Register(ModuleName, 1128, "undelegate-amount is larger than current delegation")
	ErrRedelegationOnCooldown          = sdkerrors.Register(ModuleName, 1129, "all redelegation slots are on cooldown")
	ErrMultipleRedelegationInSameBlock = sdkerrors.Register(ModuleName, 1130, "only one redelegation per delegator per block")
)
