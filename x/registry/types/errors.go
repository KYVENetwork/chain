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
	ErrUploadInterval          = sdkerrors.Register(ModuleName, 1108, "upload interval not surpassed")
	ErrInvalidBundleId        = sdkerrors.Register(ModuleName, 1109, "current bundleId %v does not match provided bundleId")
	ErrAlreadyVoted           = sdkerrors.Register(ModuleName, 1110, "already voted on proposal %v")
	ErrQuorumNotReached       = sdkerrors.Register(ModuleName, 1111, "quorum not reached")
	ErrVoterIsUploader        = sdkerrors.Register(ModuleName, 1112, "voter is uploader")
	ErrNotDesignatedUploader  = sdkerrors.Register(ModuleName, 1113, "not designated uploader")
	ErrUploaderAlreadyClaimed = sdkerrors.Register(ModuleName, 1114, "uploader role already claimed")
	ErrNotEnoughNodesOnline   = sdkerrors.Register(ModuleName, 1115, "not enough nodes online")
	ErrInvalidCommission   = sdkerrors.Register(ModuleName, 1116, "invalid commission %v")
	ErrSelfDelegation   = sdkerrors.Register(ModuleName, 1117, "self delegation not allowed")
	ErrFromHeight   = sdkerrors.Register(ModuleName, 1118, "invalid from height")
	ErrInvalidVote   = sdkerrors.Register(ModuleName, 1119, "invalid vote %v")
	ErrMaxBundleSize   = sdkerrors.Register(ModuleName, 1120, "bundle size is too high")
	ErrPoolCurrentlyUpgrading   = sdkerrors.Register(ModuleName, 1121, "pool currently upgrading")
	ErrPoolNoUpgradeScheduled   = sdkerrors.Register(ModuleName, 1122, "pool has no scheduled upgrade")
)

// delegation errors
var (
	ErrNotADelegator       = sdkerrors.Register(ModuleName, 1123, "not a delegator")
	ErrNotEnoughDelegation = sdkerrors.Register(ModuleName, 1124, "undelegate-amount is larger than current delegation")
)
