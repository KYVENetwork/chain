package keeper

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"github.com/KYVENetwork/chain/x/query/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetFullStaker(ctx sdk.Context, stakerAddress string) *types.FullStaker {
	staker, _ := k.stakerKeeper.GetValidator(ctx, stakerAddress)

	stakerMetadata := types.StakerMetadata{
		Commission:              staker.Commission.Rate,
		Moniker:                 staker.GetMoniker(),
		Website:                 staker.Description.Website,
		Identity:                staker.Description.Identity,
		SecurityContact:         staker.Description.SecurityContact,
		Details:                 staker.Description.GetDetails(),
		PendingCommissionChange: nil,
		CommissionRewards:       sdk.NewCoins(),
	}

	var poolMemberships []*types.PoolMembership

	for _, valaccount := range k.stakerKeeper.GetValaccountsFromStaker(ctx, stakerAddress) {

		pool, _ := k.poolKeeper.GetPool(ctx, valaccount.PoolId)

		accountValaddress, _ := sdk.AccAddressFromBech32(valaccount.Valaddress)
		balanceValaccount := k.bankKeeper.GetBalance(ctx, accountValaddress, globalTypes.Denom).Amount.Uint64()

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
					TotalDelegation:      k.stakerKeeper.GetTotalStakeOfPool(ctx, pool.Id),
					Status:               k.GetPoolStatus(ctx, &pool),
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
	// TODO rework query spec

	return &types.FullStaker{
		Address:                 stakerAddress,
		Metadata:                &stakerMetadata,
		SelfDelegation:          uint64(0), // TODO rework query spec
		SelfDelegationUnbonding: selfDelegationUnbonding,
		TotalDelegation:         staker.Tokens.Uint64(),
		DelegatorCount:          0, // TODO rework query spec
		Pools:                   poolMemberships,
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
