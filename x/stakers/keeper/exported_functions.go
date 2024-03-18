package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Gov
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	// Stakers
	"github.com/KYVENetwork/chain/x/stakers/types"
)

// These functions are meant to be called from external modules.
// For now this is the bundles module and the delegation module
// which need to interact with the stakers module.

// LeavePool removes a staker from a pool and emits the corresponding event.
// The staker is no longer able to participate in the given pool.
// All points the staker had in that pool are deleted.
func (k Keeper) LeavePool(ctx sdk.Context, staker string, poolId uint64) {
	k.RemoveValaccountFromPool(ctx, poolId, staker)

	_ = ctx.EventManager().EmitTypedEvent(&types.EventLeavePool{
		PoolId: poolId,
		Staker: staker,
	})
}

// GetAllStakerAddressesOfPool returns a list of all stakers
// which have currently a valaccount registered for the given pool
// and are therefore allowed to participate in that pool.
func (k Keeper) GetAllStakerAddressesOfPool(ctx sdk.Context, poolId uint64) (stakers []string) {
	for _, valaccount := range k.GetAllValaccountsOfPool(ctx, poolId) {
		stakers = append(stakers, valaccount.Staker)
	}

	return stakers
}

// GetCommission returns the commission of a staker as a parsed sdk.Dec
func (k Keeper) GetCommission(ctx sdk.Context, stakerAddress string) sdk.Dec {
	staker, _ := k.GetStaker(ctx, stakerAddress)
	return staker.Commission
}

// AssertValaccountAuthorized checks if the given `valaddress` is allowed to vote in pool
// with id `poolId` to vote in favor of `stakerAddress`.
// If the valaddress is not authorized the appropriate error is returned.
// Otherwise, it returns `nil`
func (k Keeper) AssertValaccountAuthorized(ctx sdk.Context, poolId uint64, stakerAddress string, valaddress string) error {
	valaccount, found := k.GetValaccount(ctx, poolId, stakerAddress)

	if !found {
		return types.ErrValaccountUnauthorized
	}

	if valaccount.Valaddress != valaddress {
		return types.ErrValaccountUnauthorized
	}

	return nil
}

// GetActiveStakers returns all staker-addresses that are
// currently participating in at least one pool.
func (k Keeper) GetActiveStakers(ctx sdk.Context) []string {
	return k.getAllActiveStakers(ctx)
}

// GOVERNANCE - BONDING
// The next functions are used in our custom fork of the cosmos-sdk
// which includes protocol staking into the governance.
// The behavior is exactly the same as with normal cosmos-validators.

// TotalBondedTokens returns all tokens which are currently bonded by the protocol
// I.e. the sum of all delegation of all stakers that are currently participating
// in at least one pool
func (k Keeper) TotalBondedTokens(ctx context.Context) math.Int {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	bondedTokens := math.ZeroInt()

	for _, validator := range k.getAllActiveStakers(sdkCtx) {
		delegation := int64(k.delegationKeeper.GetDelegationAmount(sdkCtx, validator))

		bondedTokens = bondedTokens.Add(math.NewInt(delegation))
	}

	return bondedTokens
}

// GetActiveValidators returns all protocol-node information which
// are needed by the governance to calculate the voting powers.
// The interface needs to correspond to github.com/cosmos/cosmos-sdk/x/gov/types/v1.ValidatorGovInfo
// But as there is no direct dependency in the cosmos-sdk-fork this value is passed as an interface{}
func (k Keeper) GetActiveValidators(ctx context.Context) (validators []interface{}) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, address := range k.getAllActiveStakers(sdkCtx) {
		delegation := int64(k.delegationKeeper.GetDelegationAmount(sdkCtx, address))

		validator := govV1Types.NewValidatorGovInfo(
			sdk.ValAddress(sdk.MustAccAddressFromBech32(address)),
			math.NewInt(delegation),
			sdk.NewDec(delegation),
			sdk.ZeroDec(),
			govV1Types.WeightedVoteOptions{},
		)

		validators = append(validators, validator)
	}

	return
}

// GetDelegations returns the address and the delegation amount of all active protocol-stakers the
// delegator as delegated to. This is used to calculate the vote weight each delegator has.
func (k Keeper) GetDelegations(ctx context.Context, delegator string) (validators []string, amounts []sdk.Dec) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, validator := range k.delegationKeeper.GetStakersByDelegator(sdkCtx, delegator) {
		if k.isActiveStaker(sdkCtx, validator) {
			validators = append(validators, validator)

			amounts = append(
				amounts,
				sdk.NewDec(int64(k.delegationKeeper.GetDelegationAmountOfDelegator(sdkCtx, validator, delegator))),
			)
		}
	}

	return
}
