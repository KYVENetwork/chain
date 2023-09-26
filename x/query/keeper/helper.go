package keeper

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetFullStaker(ctx sdk.Context, stakerAddress string) *types.FullStaker {
	staker, _ := k.stakerKeeper.GetStaker(ctx, stakerAddress)

	commissionChange, found := k.stakerKeeper.GetCommissionChangeEntryByIndex2(ctx, staker.Address)
	var commissionChangeEntry *types.CommissionChangeEntry = nil
	if found {
		commissionChangeEntry = &types.CommissionChangeEntry{
			Commission:   commissionChange.Commission,
			CreationDate: commissionChange.CreationDate,
		}
	}

	stakerMetadata := types.StakerMetadata{
		Commission:              staker.Commission,
		Moniker:                 staker.Moniker,
		Website:                 staker.Website,
		Identity:                staker.Identity,
		SecurityContact:         staker.SecurityContact,
		Details:                 staker.Details,
		PendingCommissionChange: commissionChangeEntry,
		CommissionRewards:       staker.CommissionRewards,
	}

	delegationData, _ := k.delegationKeeper.GetDelegationData(ctx, staker.Address)

	var poolMemberships []*types.PoolMembership

	for _, valaccount := range k.stakerKeeper.GetValaccountsFromStaker(ctx, staker.Address) {

		pool, _ := k.poolKeeper.GetPool(ctx, valaccount.PoolId)

		accountValaddress, _ := sdk.AccAddressFromBech32(valaccount.Valaddress)
		balanceValaccount := k.bankKeeper.GetBalance(ctx, accountValaddress, globalTypes.Denom).Amount.Uint64()

		poolMemberships = append(
			poolMemberships, &types.PoolMembership{
				Pool: &types.BasicPool{
					Id:             pool.Id,
					Name:           pool.Name,
					Runtime:        pool.Runtime,
					Logo:           pool.Logo,
					OperatingCost:  pool.OperatingCost,
					UploadInterval: pool.UploadInterval,
					// TODO(rapha): fix this
					//TotalFunds:      pool.TotalFunds,
					TotalDelegation: k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id),
					Status:          k.GetPoolStatus(ctx, &pool),
				},
				Points:     valaccount.Points,
				IsLeaving:  valaccount.IsLeaving,
				Valaddress: valaccount.Valaddress,
				Balance:    balanceValaccount,
			},
		)
	}

	// Iterate all UnbondingDelegation entries to get total delegation unbonding amount
	selfDelegationUnbonding := uint64(0)
	for _, entry := range k.delegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(ctx, stakerAddress) {
		if entry.Staker == stakerAddress {
			selfDelegationUnbonding += entry.Amount
		}
	}

	return &types.FullStaker{
		Address:                 staker.Address,
		Metadata:                &stakerMetadata,
		SelfDelegation:          k.delegationKeeper.GetDelegationAmountOfDelegator(ctx, stakerAddress, stakerAddress),
		SelfDelegationUnbonding: selfDelegationUnbonding,
		TotalDelegation:         k.delegationKeeper.GetDelegationAmount(ctx, staker.Address),
		DelegatorCount:          delegationData.DelegatorCount,
		Pools:                   poolMemberships,
	}
}

func (k Keeper) GetPoolStatus(ctx sdk.Context, pool *pooltypes.Pool) pooltypes.PoolStatus {
	totalDelegation := k.delegationKeeper.GetDelegationOfPool(ctx, pool.Id)

	var poolStatus pooltypes.PoolStatus

	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		poolStatus = pooltypes.POOL_STATUS_UPGRADING
	} else if pool.Disabled {
		poolStatus = pooltypes.POOL_STATUS_DISABLED
	} else if totalDelegation < pool.MinDelegation {
		poolStatus = pooltypes.POOL_STATUS_NOT_ENOUGH_DELEGATION
		// TODO(rapha): fix this
		//} else if pool.TotalFunds == 0 {
		//	poolStatus = pooltypes.POOL_STATUS_NO_FUNDS
	} else {
		poolStatus = pooltypes.POOL_STATUS_ACTIVE
	}

	return poolStatus
}
