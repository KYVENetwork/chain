package types

import (
	"github.com/KYVENetwork/chain/util"
	stakersTypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	// Bundles
	bundlesTypes "github.com/KYVENetwork/chain/x/bundles/types"
	// Delegation
	delegationTypes "github.com/KYVENetwork/chain/x/delegation/types"
	// Global
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// Pool
	poolTypes "github.com/KYVENetwork/chain/x/pool/types"
)

type BundlesKeeper interface {
	util.BundlesKeeper

	GetBundleProposal(sdk.Context, uint64) (bundlesTypes.BundleProposal, bool)
	GetFinalizedBundle(sdk.Context, uint64, uint64) (bundlesTypes.FinalizedBundle, bool)
	GetFinalizedBundleByHeight(sdk.Context, uint64, uint64) (bundlesTypes.FinalizedBundle, bool)
	GetPaginatedFinalizedBundleQuery(sdk.Context, *query.PageRequest, uint64) ([]bundlesTypes.FinalizedBundle, *query.PageResponse, error)
	GetParams(sdk.Context) bundlesTypes.Params
	GetVoteDistribution(sdk.Context, uint64) bundlesTypes.VoteDistribution
}

type DelegationKeeper interface {
	util.DelegationKeeper

	GetAllUnbondingDelegationQueueEntriesOfDelegator(sdk.Context, string) []delegationTypes.UndelegationQueueEntry
	GetDelegationData(sdk.Context, string) (delegationTypes.DelegationData, bool)
	GetParams(sdk.Context) delegationTypes.Params
	GetUndelegationQueueEntry(sdk.Context, uint64) (delegationTypes.UndelegationQueueEntry, bool)
}

type GlobalKeeper interface {
	util.GlobalKeeper

	GetParams(sdk.Context) globalTypes.Params
}

type PoolKeeper interface {
	util.PoolKeeper

	GetPaginatedPoolsQuery(sdk.Context, *query.PageRequest, string, string, bool, uint32) ([]poolTypes.Pool, *query.PageResponse, error)
	GetAllPools(sdk.Context) []poolTypes.Pool
	GetPool(sdk.Context, uint64) (poolTypes.Pool, bool)
	GetPoolWithError(sdk.Context, uint64) (poolTypes.Pool, error)
}

type StakersKeeper interface {
	util.StakersKeeper

	GetAllValaccountsOfPool(sdk.Context, uint64) []*stakersTypes.Valaccount
	GetCommissionChangeEntryByIndex2(sdk.Context, string) (stakersTypes.CommissionChangeEntry, bool)
	GetParams(sdk.Context) stakersTypes.Params
	GetStaker(sdk.Context, string) (stakersTypes.Staker, bool)
	GetValaccountsFromStaker(sdk.Context, string) []*stakersTypes.Valaccount
}
