package bundles

import (
	"github.com/KYVENetwork/chain/util"
	bundlesKeeper "github.com/KYVENetwork/chain/x/bundles/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// Mint
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	// Pool
	"github.com/KYVENetwork/chain/x/pool/keeper"
	// Team
	teamKeeper "github.com/KYVENetwork/chain/x/team/keeper"
	// Upgrade
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
)

func SplitInflation(ctx sdk.Context, k bundlesKeeper.Keeper, bk bankKeeper.Keeper, mk mintKeeper.Keeper, pk keeper.Keeper, tk teamKeeper.Keeper, uk upgradeKeeper.Keeper) {
	minter := mk.GetMinter(ctx)
	params := mk.GetParams(ctx)

	// get total inflation rewards for current block
	blockProvision := minter.BlockProvision(params).Amount.Int64()

	// calculate the remaining block provision for chain and protocol after x/team took its share
	remainingBlockProvision := blockProvision - tk.GetTeamBlockProvision(ctx)

	// calculate block provision for protocol based on protocol inflation share
	protocolBlockProvision := sdk.NewDec(remainingBlockProvision).Mul(pk.GetProtocolInflationShare(ctx)).TruncateInt64()

	if protocolBlockProvision == 0 {
		return
	}

	// track actual distributed block provision for protocol
	distributed := uint64(0)

	// calculate total operating cost of pools to get each pool's reward share
	totalOperatingCost := uint64(0)

	for _, pool := range pk.GetAllPools(ctx) {
		// only include active pools
		if err := k.AssertPoolCanRun(ctx, pool.Id); err == nil {
			totalOperatingCost += pool.OperatingCost
		}
	}

	// if the total operating cost is zero all rewards go the chain
	if totalOperatingCost == 0 {
		return
	}

	for _, pool := range pk.GetAllPools(ctx) {
		// only include active pools
		if err := k.AssertPoolCanRun(ctx, pool.Id); err == nil {
			// calculate pool share based of operating cost
			amount := uint64(sdk.NewDec(int64(pool.OperatingCost)).Quo(sdk.NewDec(int64(totalOperatingCost))).Mul(sdk.NewDec(protocolBlockProvision)).TruncateInt64())

			// transfer funds to pool account
			if err := util.TransferFromModuleToAddress(bk, ctx, authTypes.FeeCollectorName, pool.GetPoolAccount().String(), amount); err != nil {
				util.PanicHalt(uk, ctx, err.Error())
			}

			// track transferred $KYVE to protocol
			distributed += amount
		}
	}

	// calculate if a remainder is left
	remainder := uint64(protocolBlockProvision) - distributed

	if remainder > 0 {
		// find an active pool
		for _, pool := range pk.GetAllPools(ctx) {
			if err := k.AssertPoolCanRun(ctx, pool.Id); err != nil {
				// add remainder to first active pool we find
				if err := util.TransferFromModuleToAddress(bk, ctx, authTypes.FeeCollectorName, pool.GetPoolAccount().String(), remainder); err != nil {
					util.PanicHalt(uk, ctx, err.Error())
				}
			}
		}
	}

	// remainder gets transferred to chain
	pk.Logger(ctx).Info("split portion of minted coins to protocol", "amount", protocolBlockProvision)
}
