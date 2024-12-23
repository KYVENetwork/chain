package keeper

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetFullStaker(ctx sdk.Context, stakerAddress string) *types.FullStaker {
	validator, _ := k.stakerKeeper.GetValidator(ctx, stakerAddress)

	var poolMemberships []*types.PoolMembership
	totalPoolStake := uint64(0)

	for _, poolAccount := range k.stakerKeeper.GetPoolAccountsFromStaker(ctx, stakerAddress) {
		pool, _ := k.poolKeeper.GetPool(ctx, poolAccount.PoolId)

		accountValaddress, _ := sdk.AccAddressFromBech32(poolAccount.PoolAddress)
		balanceValaccount := k.bankKeeper.GetBalance(ctx, accountValaddress, globalTypes.Denom).Amount.Uint64()

		commissionChange, found := k.stakerKeeper.GetCommissionChangeEntryByIndex2(ctx, stakerAddress, poolAccount.PoolId)
		var commissionChangeEntry *types.CommissionChangeEntry = nil
		if found {
			commissionChangeEntry = &types.CommissionChangeEntry{
				Commission:   commissionChange.Commission,
				CreationDate: commissionChange.CreationDate,
			}
		}

		stakeFractionChange, found := k.stakerKeeper.GetStakeFractionChangeEntryByIndex2(ctx, stakerAddress, poolAccount.PoolId)
		var stakeFractionChangeEntry *types.StakeFractionChangeEntry = nil
		if found {
			stakeFractionChangeEntry = &types.StakeFractionChangeEntry{
				StakeFraction: stakeFractionChange.StakeFraction,
				CreationDate:  stakeFractionChange.CreationDate,
			}
		}

		poolStake := k.stakerKeeper.GetValidatorPoolStake(ctx, stakerAddress, pool.Id)
		totalPoolStake += poolStake

		poolMemberships = append(
			poolMemberships, &types.PoolMembership{
				Pool: &types.BasicPool{
					Id:                   pool.Id,
					Name:                 pool.Name,
					Runtime:              pool.Runtime,
					Logo:                 pool.Logo,
					InflationShareWeight: pool.InflationShareWeight,
					UploadInterval:       pool.UploadInterval,
					TotalFunds:           k.fundersKeeper.GetTotalActiveFunding(ctx, pool.Id),
					TotalStake:           k.stakerKeeper.GetTotalStakeOfPool(ctx, pool.Id),
					Status:               k.GetPoolStatus(ctx, &pool),
				},
				Points:                     poolAccount.Points,
				IsLeaving:                  poolAccount.IsLeaving,
				PoolAddress:                poolAccount.PoolAddress,
				Balance:                    balanceValaccount,
				Commission:                 poolAccount.Commission,
				PendingCommissionChange:    commissionChangeEntry,
				StakeFraction:              poolAccount.StakeFraction,
				PendingStakeFractionChange: stakeFractionChangeEntry,
				PoolStake:                  poolStake,
			},
		)
	}

	return &types.FullStaker{
		Address:        stakerAddress,
		Validator:      &validator,
		TotalPoolStake: totalPoolStake,
		Pools:          poolMemberships,
	}
}

func (k Keeper) GetPoolStatus(ctx sdk.Context, pool *pooltypes.Pool) pooltypes.PoolStatus {
	var poolStatus pooltypes.PoolStatus

	poolStatus = pooltypes.POOL_STATUS_ACTIVE
	if pool.UpgradePlan.ScheduledAt > 0 && uint64(ctx.BlockTime().Unix()) >= pool.UpgradePlan.ScheduledAt {
		poolStatus = pooltypes.POOL_STATUS_UPGRADING
	} else if pool.Disabled {
		poolStatus = pooltypes.POOL_STATUS_DISABLED
	} else if pool.EndKey != "" && pool.EndKey == pool.CurrentKey {
		poolStatus = pooltypes.POOL_STATUS_END_KEY_REACHED
	} else if k.stakerKeeper.IsVotingPowerTooHigh(ctx, pool.Id) {
		poolStatus = pooltypes.POOL_STATUS_VOTING_POWER_TOO_HIGH
	} else if k.stakerKeeper.GetTotalStakeOfPool(ctx, pool.Id) < pool.MinDelegation {
		poolStatus = pooltypes.POOL_STATUS_NOT_ENOUGH_DELEGATION
	} else if k.fundersKeeper.GetTotalActiveFunding(ctx, pool.Id).IsZero() {
		poolStatus = pooltypes.POOL_STATUS_NO_FUNDS
	}

	return poolStatus
}
