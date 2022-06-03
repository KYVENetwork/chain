package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func emitBasicEvent(ctx sdk.Context, eventName string, attrs ...sdk.Attribute) {
	fullAttrs := []sdk.Attribute{}
	fullAttrs = append(fullAttrs, sdk.NewAttribute(EventName, eventName))
	fullAttrs = append(fullAttrs, sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			append(fullAttrs, attrs...)...,
		),
	)
}

func emitBundleEvent(ctx sdk.Context, pool *Pool, status string, bundleReward uint64) {
	emitBasicEvent(ctx, ProposalEventKey,
		sdk.NewAttribute(EventPoolId, strconv.FormatUint(pool.Id, 10)),
		sdk.NewAttribute(ProposalEventBundleId, pool.BundleProposal.BundleId),
		sdk.NewAttribute(ProposalEventByteSize, strconv.FormatUint(pool.BundleProposal.ByteSize, 10)),
		sdk.NewAttribute(ProposalEventUploader, pool.BundleProposal.Uploader),
		sdk.NewAttribute(ProposalEventNextUploader, pool.BundleProposal.NextUploader),
		sdk.NewAttribute(ProposalEventReward, strconv.FormatUint(bundleReward, 10)),
		sdk.NewAttribute(ProposalEventValid, strconv.FormatUint(uint64(len(pool.BundleProposal.VotersValid)), 10)),
		sdk.NewAttribute(ProposalEventInvalid, strconv.FormatUint(uint64(len(pool.BundleProposal.VotersInvalid)), 10)),
		sdk.NewAttribute(ProposalEventFromHeight, strconv.FormatUint(pool.BundleProposal.FromHeight, 10)),
		sdk.NewAttribute(ProposalEventToHeight, strconv.FormatUint(pool.BundleProposal.ToHeight, 10)),
		sdk.NewAttribute(ProposalEventStatus, status),
	)
}

// emitMoneyPoolTransferEvent
// Event Schema which can be used for (un)staking and (de)funding
func emitMoneyPoolTransferEvent(ctx sdk.Context, eventName string, poolId uint64, account string, amount uint64) {
	emitBasicEvent(ctx, eventName,
		sdk.NewAttribute(EventPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(EventCreator, account),
		sdk.NewAttribute(EventAmount, strconv.FormatUint(amount, 10)),
	)
}

// emitMoneyPoolStakerTransferEvent
// Event Schema which can be used for (un)delegating
func emitMoneyPoolStakerTransferEvent(ctx sdk.Context, eventName string, poolId uint64, account string, staker string, amount uint64) {
	emitBasicEvent(ctx, eventName,
		sdk.NewAttribute(EventPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute("Staker", staker),
		sdk.NewAttribute("Creator", account),
		sdk.NewAttribute(EventAmount, strconv.FormatUint(amount, 10)),
	)
}

// FUNDING

func EmitFundPoolEvent(ctx sdk.Context, msg *MsgFundPool) {
	emitMoneyPoolTransferEvent(ctx, "Funded", msg.Id, msg.Creator, msg.Amount)
}

func EmitDefundPoolEvent(ctx sdk.Context, poolId uint64, account string, amount uint64) {
	emitMoneyPoolTransferEvent(ctx, "Defunded", poolId, account, amount)
}

// STAKING

func EmitStakeEvent(ctx sdk.Context, creator string, poolId uint64, amount uint64) {
	emitMoneyPoolTransferEvent(ctx, "Staked", poolId, creator, amount)
}

func EmitUnstakeEvent(ctx sdk.Context, poolId uint64, stakerCreator string, amount uint64) {
	emitMoneyPoolTransferEvent(ctx, "Unstaked", poolId, stakerCreator, amount)
}

func EmitUpdateMetadata(ctx sdk.Context, creator string, poolId uint64, commission string, moniker string, website string, logo string) {
	emitBasicEvent(ctx, UpdateMetadataEventKey,
		sdk.NewAttribute(EventCreator, creator),
		sdk.NewAttribute(EventPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(UpdateMetadataCommission, commission),
		sdk.NewAttribute(UpdateMetadataMoniker, moniker),
		sdk.NewAttribute(UpdateMetadataWebsite, website),
		sdk.NewAttribute(UpdateMetadataLogo, logo),
	)
}

// DELEGATION

func EmitDelegateEvent(ctx sdk.Context, poolId uint64, creator string, staker string, amount uint64) {
	emitMoneyPoolStakerTransferEvent(ctx, "Delegated", poolId, creator, staker, amount)
}

func EmitUndelegateEvent(ctx sdk.Context, poolId uint64, creator string, staker string, amount uint64) {
	emitMoneyPoolStakerTransferEvent(ctx, "Undelegated", poolId, creator, staker, amount)
}

// BUNDLE PROPOSAL

func EmitBundleInvalidEvent(ctx sdk.Context, pool *Pool) {
	emitBundleEvent(ctx, pool, "Invalid", 0)
}

func EmitBundleValidEvent(ctx sdk.Context, pool *Pool, reward uint64) {
	emitBundleEvent(ctx, pool, "Valid", reward)
}

func EmitBundleDroppedInsufficientFundsEvent(ctx sdk.Context, pool *Pool) {
	emitBundleEvent(ctx, pool, "Dropped: not enough funds", 0)
}

func EmitBundleDroppedQuorumNotReachedEvent(ctx sdk.Context, pool *Pool) {
	emitBundleEvent(ctx, pool, "Dropped: quorum not reached", 0)
}

func EmitBundleTimeoutEvent(ctx sdk.Context, pool *Pool) {
	emitBundleEvent(ctx, pool, "Timeout", 0)
}

func EmitEmptyBundleEvent(ctx sdk.Context, pool *Pool) {
	emitBundleEvent(ctx, pool, "Empty", 0)
}

func EmitBundleVoteEvent(ctx sdk.Context, pool *Pool, msg *MsgVoteProposal) {
	emitBasicEvent(ctx, VoteEventKey,
		sdk.NewAttribute(EventPoolId, strconv.FormatUint(pool.Id, 10)),
		sdk.NewAttribute(EventCreator, msg.Creator),
		sdk.NewAttribute(VoteEventBundleId, msg.BundleId),
		sdk.NewAttribute(VoteEventVote, strconv.FormatUint(uint64(msg.Vote), 10)),
	)
}

func EmitSlashEvent(ctx sdk.Context, poolId uint64, slashedAccount string, slashAmount uint64) {
	emitBasicEvent(ctx, SlashEventKey,
		sdk.NewAttribute(EventPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(SlashAccount, slashedAccount),
		sdk.NewAttribute(EventAmount, strconv.FormatUint(slashAmount, 10)),
	)
}
