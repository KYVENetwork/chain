package v0_5_0

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

func createUnbondingParameters(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	registryKeeper.ParamStore().Set(ctx, types.KeyUnbondingStakingTime, types.DefaultUnbondingStakingTime)

	registryKeeper.ParamStore().Set(ctx, types.KeyUnbondingDelegationTime, types.DefaultUnbondingDelegationTime)
}

func createProposalIndex(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	fmt.Printf("%sCreating thrid proposal index\n", MigrationLoggerPrefix)

	// Set all delegators again to create the index
	proposals := registryKeeper.GetAllProposal(ctx)
	for index, proposal := range proposals {

		registryKeeper.SetProposal(ctx, proposal)

		if index%1000 == 0 {
			fmt.Printf("%sProposals processed: %d\n", MigrationLoggerPrefix, index)
		}
	}

	fmt.Printf("%sFinished index creation\n", MigrationLoggerPrefix)
}

func migrateIBCDenoms(ctx sdk.Context, transferKeeper *ibctransferkeeper.Keeper) {
	var newTraces []ibctransfertypes.DenomTrace

	transferKeeper.IterateDenomTraces(ctx,
		func(dt ibctransfertypes.DenomTrace) bool {
			newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())

			if err := newTrace.Validate(); err == nil && !equalTraces(newTrace, dt) {
				newTraces = append(newTraces, newTrace)
			}

			return false
		},
	)

	for _, nt := range newTraces {
		transferKeeper.SetDenomTrace(ctx, nt)
	}
}

func migratePools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	registryKeeper.ParamStore().Set(ctx, types.KeyStorageCost, uint64(50))

	for _, pool := range registryKeeper.GetAllPool(ctx) {
		// deprecate pool versions
		pool.Versions = ""

		// set 2.5 $KYVE as operating cost
		pool.OperatingCost = 2_500_000_000

		// schedule upgrades for each runtime
		switch pool.Runtime {
		case "@kyve/evm":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "1.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/evm/releases/download/v1.2.0/kyve-linux.zip?checksum=a64febdda593950a222c3b884bf4220832f2906cf7923570fd1e51de043ca09e\",\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.2.0/kyve-macos.zip?checksum=a876fbb41e3fac985880062c5fe824e96d7e6ca10380d262161d460a4bf133e2\"}",
			}
		case "@kyve/stacks":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.2.0/kyve-linux.zip?checksum=fecd7fa526f3ed5836bd0c848f3e445cf08e48797d3b587d4426f61fe4f988c4\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.2.0/kyve-macos.zip?checksum=6c4847c29beaa863daedb80534d2f04eea95840e4b32acafce2602052391dc55\"}",
			}
		case "@kyve/bitcoin":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.2.0/kyve-linux.zip?checksum=36447e50ea9eaa77e81618f45d8c40a57ae73d0a5b6abbe41e3f80879fc5a52f\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.2.0/kyve-macos.zip?checksum=5d6c73021863012a2d22e17b2aed0227776677d1ac23e79786f46a1b3ab72e05\"}",
			}
		case "@kyve/solana":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.2.0/kyve-linux.zip?checksum=2f4d13d249d890d38beaf1694d91a1fd871798dba8a954d4a96177899b51237c\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.2.0/kyve-macos.zip?checksum=8d25178917758a3fd22b48bb1368a2df801fb731187adaf6f8be0d35c8134558\"}",
			}
		case "@kyve/zilliqa":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.2.0/kyve-linux.zip?checksum=3c0a06276ce0162f7d6f0b91f00c0ddcbdf605f952f0c1a7567519aadb6dcf47\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.2.0/kyve-macos.zip?checksum=d5c5343c76e0faac2487e0b4bd5e8f47f42816a5f9813196f8f803a1695d6fe4\"}",
			}
		case "@kyve/near":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.2.0/kyve-linux.zip?checksum=57da2ca0fbcac9ddb61b1785dbc82ae1bfe27814158b0934ef0735ebc5e23033\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.2.0/kyve-macos.zip?checksum=9a383316940c81cf6f3c5b04d03427d97eacff7d326be5ec24a33b461287fb89\"}",
			}
		case "@kyve/celo":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.2.0/kyve-linux.zip?checksum=dd75d1b0d98ca3769befaa9549a10aafe91740844479e6db31e59f63bd64a5b8\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.2.0/kyve-macos.zip?checksum=a7c93222950bc5c6dadb0efa604ab02af452ad23860cf0d70901cacd07115fdb\"}",
			}
		case "@kyve/cosmos":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.2.0/kyve-linux.zip?checksum=4932493fbf5bf4896d15fe163babb9a21df55bfd855049b6183d0666de0aad9a\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.2.0/kyve-macos.zip?checksum=3577edd34a022544402161eafff6773ef8f2c779a3180acd22730964755a0eab\"}",
			}
		case "@kyve/substrate":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.1.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/substrate/releases/download/v0.1.0/kyve-linux.zip?checksum=1071e16e9563eebd76ad2910006795feeff0008e2fe30ed0715e58570cd89875\",\"macos\":\"https://github.com/kyve-org/substrate/releases/download/v0.1.0/kyve-macos.zip?checksum=0135f2bbd4cae57988ae383560d31f323d778719150144bc156ba221570ac403\"}",
			}
		default:
			pool.UpgradePlan = &types.UpgradePlan{}
		}

		// add pool upgrade info
		pool.UpgradePlan.ScheduledAt = uint64(ctx.BlockTime().Unix())
		pool.UpgradePlan.Duration = 1800 // 30min

		// save changes
		registryKeeper.SetPool(ctx, pool)
	}
}

func migrateProposals(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {

	fmt.Printf("%sMigration Proposals to new key system\n", MigrationLoggerPrefix)

	for _, pool := range registryKeeper.GetAllPool(ctx) {
		proposalPrefixBuilder := types.KeyPrefixBuilder{Key: types.ProposalKeyPrefixIndex3}.AInt(pool.Id)
		store := prefix.NewStore(ctx.KVStore(registryKeeper.StoreKey()), proposalPrefixBuilder.Key)
		iterator := sdk.KVStorePrefixIterator(store, []byte{})

		defer iterator.Close()
		id := uint64(0)

		var lastProposal types.Proposal

		for ; iterator.Valid(); iterator.Next() {
			bundleId := string(iterator.Value())
			proposal, _ := registryKeeper.GetProposal(ctx, bundleId)
			proposal.Id = id
			id += 1
			proposal.Key = strconv.FormatUint(proposal.ToHeight-1, 10)
			registryKeeper.SetProposal(ctx, proposal)

			if id%100 == 0 {
				fmt.Printf("%sPool %d : Proposals processed: %d\n", MigrationLoggerPrefix, pool.Id, id)
			}

			lastProposal = proposal
		}

		pool.TotalBundles = id

		// drop current bundle
		pool.BundleProposal = &types.BundleProposal{
			NextUploader: pool.BundleProposal.NextUploader,
			CreatedAt:    uint64(ctx.BlockHeight()),
		}

		// reset height
		pool.CurrentHeight = lastProposal.ToHeight
		// migrate height to custom keys
		pool.StartKey = strconv.FormatUint(pool.CurrentHeight, 10)

		registryKeeper.SetPool(ctx, pool)
	}

	fmt.Printf("%sFinished proposal migration\n", MigrationLoggerPrefix)
}

func updateGovParams(ctx sdk.Context, govKeeper *govkeeper.Keeper) {
	govKeeper.SetDepositParams(ctx, govtypes.DepositParams{
		// 20,000 $KYVE
		MinDeposit: sdk.NewCoins(sdk.NewInt64Coin("tkyve", 20_000_000_000_000)),
		// 5 minutes
		MaxDepositPeriod: time.Minute * 5,
		// 100,000 $KYVE
		MinExpeditedDeposit: sdk.NewCoins(sdk.NewInt64Coin("tkyve", 100_000_000_000_000)),
	})

	govKeeper.SetTallyParams(ctx, govtypes.TallyParams{
		// 0.01 - 1%
		Quorum: sdk.NewDec(1).Quo(sdk.NewDec(100)),
		// 0.5 - 50%
		Threshold: sdk.NewDecWithPrec(5, 1),
		// 0.334 - 33.4%
		VetoThreshold: sdk.NewDecWithPrec(334, 3),
		// 0.667 - 66.7%
		ExpeditedThreshold: sdk.NewDecWithPrec(667, 3),
	})

	govKeeper.SetVotingParams(ctx, govtypes.VotingParams{
		VotingPeriod: time.Minute * 60 * 24,
		ProposalVotingPeriods: []govtypes.ProposalVotingPeriod{
			{
				ProposalType: "kyve.registry.v1beta1.CreatePoolProposal",
				VotingPeriod: time.Minute * 60 * 2,
			},
		},
		ExpeditedVotingPeriod: time.Minute * 30,
	})
}

func CreateUpgradeHandler(
	govKeeper *govkeeper.Keeper,
	registryKeeper *registrykeeper.Keeper,
	transferKeeper *ibctransferkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		createUnbondingParameters(registryKeeper, ctx)

		createProposalIndex(registryKeeper, ctx)

		migrateIBCDenoms(ctx, transferKeeper)

		updateGovParams(ctx, govKeeper)

		migratePools(registryKeeper, ctx)

		migrateProposals(registryKeeper, ctx)

		// Return.
		return vm, nil
	}
}

func equalTraces(dtA, dtB ibctransfertypes.DenomTrace) bool {
	return dtA.BaseDenom == dtB.BaseDenom && dtA.Path == dtB.Path
}
