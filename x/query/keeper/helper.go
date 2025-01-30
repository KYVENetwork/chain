package keeper

import (
	"cosmossdk.io/math"

	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) GetFullStaker(ctx sdk.Context, stakerAddress string) (*types.FullStaker, error) {
	valAddress, err := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(stakerAddress))
	if err != nil {
		return nil, err
	}
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddress)
	if err != nil {
		return nil, err
	}

	validatorAccAddress, _ := sdk.AccAddressFromBech32(stakerAddress)
	validatorUnbonding := math.ZeroInt()
	err = k.stakingKeeper.IterateDelegatorUnbondingDelegations(ctx, validatorAccAddress, func(ubd stakingTypes.UnbondingDelegation) bool {
		if ubd.ValidatorAddress == validator.OperatorAddress {
			for _, entry := range ubd.Entries {
				validatorUnbonding = validatorUnbonding.Add(entry.Balance)
			}
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	validatorDelegations, err := k.stakingKeeper.GetValidatorDelegations(ctx, valAddress)
	if err != nil {
		return nil, err
	}

	validatorSelfDelegation := uint64(0)

	if d, err := k.stakingKeeper.GetDelegation(ctx, validatorAccAddress, valAddress); err == nil {
		validatorSelfDelegation = uint64(validator.TokensFromSharesTruncated(d.Shares).TruncateInt64())
	}

	validatorCommissionRewards, err := k.distrkeeper.GetValidatorAccumulatedCommission(ctx, valAddress)
	if err != nil {
		return nil, err
	}

	var poolMemberships []*types.PoolMembership
	validatorTotalPoolStake := uint64(0)

	for _, poolAccount := range k.stakerKeeper.GetPoolAccountsFromStaker(ctx, stakerAddress) {
		pool, _ := k.poolKeeper.GetPool(ctx, poolAccount.PoolId)

		poolAddressAccount, _ := sdk.AccAddressFromBech32(poolAccount.PoolAddress)
		balancePoolAddress := k.bankKeeper.GetBalance(ctx, poolAddressAccount, globalTypes.Denom).Amount.Uint64()

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
		validatorTotalPoolStake += poolStake

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
				Balance:                    balancePoolAddress,
				Commission:                 poolAccount.Commission,
				PendingCommissionChange:    commissionChangeEntry,
				StakeFraction:              poolAccount.StakeFraction,
				PendingStakeFractionChange: stakeFractionChangeEntry,
				PoolStake:                  poolStake,
			},
		)
	}

	return &types.FullStaker{
		Address:                    stakerAddress,
		Validator:                  &validator,
		ValidatorDelegators:        uint64(len(validatorDelegations)),
		ValidatorUnbonding:         validatorUnbonding.Uint64(),
		ValidatorSelfDelegation:    validatorSelfDelegation,
		ValidatorTotalPoolStake:    validatorTotalPoolStake,
		ValidatorCommissionRewards: util.TruncateDecCoins(validatorCommissionRewards.Commission),
		Pools:                      poolMemberships,
	}, nil
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
