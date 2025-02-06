package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"
	fundersTypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - valid bundles

* Produce a valid bundle with multiple validators and no foreign delegations
* Produce a valid bundle with multiple validators and foreign delegations
* Produce a valid bundle with multiple validators and foreign delegation although some did not vote at all
* Produce a valid bundle with multiple validators and foreign delegation although some voted abstain
* Produce a valid bundle with multiple validators and foreign delegation although some voted invalid
* Produce a valid bundle with multiple validators and foreign delegation although some voted invalid with maximum voting power
* Produce a valid bundle with multiple validators and no foreign delegations and another storage provider
* Produce a valid bundle with multiple validators, multiple coins and no foreign delegations
* Produce a valid bundle with multiple validators, multiple coins which are not enough for the storage reward and no foreign delegations
* Produce a valid bundle with multiple validator, multiple coins and no foreign delegations where one coin is removed from the whitelist
* Produce a valid bundle with multiple validators, multiple coins and real world values

*/

var _ = Describe("valid bundles", Ordered, func() {
	var s *i.KeeperTestSuite
	var initialBalanceStaker0, initialBalancePoolAddress0, initialBalanceStaker1, initialBalancePoolAddress1, initialBalanceStaker2, initialBalancePoolAddress2 sdk.Coins

	amountPerBundle := int64(10_000)
	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		initialBalanceStaker0 = s.GetCoinsFromAddress(i.STAKER_0)
		initialBalancePoolAddress0 = s.GetCoinsFromAddress(i.POOL_ADDRESS_0_A)

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		initialBalanceStaker2 = s.GetCoinsFromAddress(i.STAKER_2)
		initialBalancePoolAddress2 = s.GetCoinsFromAddress(i.POOL_ADDRESS_2_A)

		// create clean pool for every test case
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: math.LegacyNewDec(10_000),
			MinDelegation:        0 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    1,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		// create funders
		s.RunTxFundersSuccess(&fundersTypes.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "Alice",
		})

		// set storage cost to 0.5
		bundleParams := s.App().BundlesKeeper.GetParams(s.Ctx())
		bundleParams.StorageCosts = append(bundleParams.StorageCosts, bundletypes.StorageCost{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.5")})
		s.App().BundlesKeeper.SetParams(s.Ctx(), bundleParams)

		// set funders params
		s.App().FundersKeeper.SetParams(s.Ctx(), fundersTypes.NewParams([]*fundersTypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(3),
			},
		}, 0))

		s.RunTxPoolSuccess(&fundersTypes.MsgFundPool{
			Creator:          i.ALICE,
			Amounts:          i.ACoins(100 * i.T_KYVE),
			AmountsPerBundle: i.ACoins(amountPerBundle),
		})

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.POOL_ADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		initialBalanceStaker0 = s.GetCoinsFromAddress(i.STAKER_0)
		initialBalancePoolAddress0 = s.GetCoinsFromAddress(i.POOL_ADDRESS_0_A)

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a valid bundle with multiple validators and no foreign delegations", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(200 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert uploader self delegation rewards
		// (total_bundle_payout - treasury_reward - storage_cost) * (1 - commission)
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(8865).String()))
		// assert commission rewards
		// (total_bundle_payout - treasury_reward - storage_cost) * commission + storage_cost
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1035).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators and foreign delegations", func() {
		// ARRANGE
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.BOB,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(800 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(800 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1035).String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (1/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(2216).String()))
		// assert delegator delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (3/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.ALICE).String()).To(Equal(i.ACoins(6648).String()))

		// check voter rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.BOB)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators and foreign delegation although some did not vote at all", func() {
		// ARRANGE
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.BOB,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.CHARLIE,
			util.MustValaddressFromOperatorAddress(i.STAKER_2),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(800 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(1200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(poolAccountVoter.Points).To(Equal(uint64(1)))

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1035).String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (1/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(2216).String()))
		// assert delegator delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (3/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.ALICE).String()).To(Equal(i.ACoins(6648).String()))

		// check voter rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators and foreign delegation although some voted abstain", func() {
		// ARRANGE
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.BOB,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.CHARLIE,
			util.MustValaddressFromOperatorAddress(i.STAKER_2),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(800 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(1200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1035).String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (1/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(2216).String()))
		// assert delegator delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (3/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.ALICE).String()).To(Equal(i.ACoins(6648).String()))

		// check voter rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators and foreign delegation although some voted invalid", func() {
		// ARRANGE
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.BOB,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.CHARLIE,
			util.MustValaddressFromOperatorAddress(i.STAKER_2),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		initialBalanceStaker2 = s.GetCoinsFromAddress(i.STAKER_2)
		initialBalancePoolAddress2 = s.GetCoinsFromAddress(i.POOL_ADDRESS_2_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(800 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(1200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// calculate voter slashes
		fraction := s.App().StakersKeeper.GetVoteSlash(s.Ctx())
		slashAmountVoter := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())
		slashAmountDelegator := uint64(math.LegacyNewDec(int64(300 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)).To(Equal(100*i.KYVE - slashAmountVoter))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(Equal(300*i.KYVE - slashAmountDelegator))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(800 * i.KYVE))

		// check voter status
		_, voterActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(voterActive).To(BeFalse())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(i.POOL_ADDRESS_2_A)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress2))

		balanceVoter := s.GetCoinsFromAddress(i.STAKER_2)
		Expect(balanceVoter).To(Equal(initialBalanceStaker2))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1035).String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (1/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(2216).String()))
		// assert delegator delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1) * (3/4)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.ALICE).String()).To(Equal(i.ACoins(6648).String()))

		// check voter rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators and foreign delegation although some voted invalid with maximum voting power", func() {
		// ARRANGE
		s.SetMaxVotingPower("0.4")

		s.CreateValidator(i.STAKER_2, "Staker-2", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_2,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_2_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// delegate to staker 2 so he runs into the max voting power limit
		s.RunTxSuccess(stakingTypes.NewMsgDelegate(
			i.CHARLIE,
			util.MustValaddressFromOperatorAddress(i.STAKER_2),
			sdk.NewInt64Coin(globalTypes.Denom, int64(300*i.KYVE)),
		))

		s.CreateValidator(i.STAKER_3, "Staker-3", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       i.STAKER_3,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_3_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.POOL_ADDRESS_0_A,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_2_A,
			Staker:    i.STAKER_2,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_INVALID,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_3_A,
			Staker:    i.STAKER_3,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker2 = s.GetCoinsFromAddress(i.STAKER_2)
		initialBalancePoolAddress2 = s.GetCoinsFromAddress(i.POOL_ADDRESS_2_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		// TODO: why is staker 2 selected as next uploader?
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_2_A,
			Staker:        i.STAKER_2,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(300 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(500*i.KYVE + 1))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_2))
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_3))
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_2))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// calculate voter slashes (due to maximum vote power only 200 kyve out of 400 where at risk for slashing)
		fraction := s.App().StakersKeeper.GetVoteSlash(s.Ctx()).Mul(math.LegacyMustNewDecFromStr("0.5")) // 200 / 400
		slashAmountVoter := uint64(math.LegacyNewDec(int64(100 * i.KYVE)).Mul(fraction).TruncateInt64())
		slashAmountDelegator := uint64(math.LegacyNewDec(int64(300 * i.KYVE)).Mul(fraction).TruncateInt64())

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.STAKER_2)).To(Equal(100*i.KYVE - slashAmountVoter))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(Equal(300*i.KYVE - slashAmountDelegator))

		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal(300 * i.KYVE))

		// check voter status
		_, voterActive := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_2, 0)
		Expect(voterActive).To(BeFalse())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(i.POOL_ADDRESS_2_A)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress2))

		balanceVoter := s.GetCoinsFromAddress(i.STAKER_2)
		Expect(balanceVoter).To(Equal(initialBalanceStaker2))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * 0.1 + (100 * 0.5)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1035).String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.5)) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(8865).String()))

		// check voter rewards
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_2, i.CHARLIE)).To(BeEmpty())

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators and no foreign delegations and another storage provider", func() {
		// ARRANGE
		storageProviderId := uint32(2)

		params := s.App().BundlesKeeper.GetParams(s.Ctx())
		params.StorageCosts = append(params.StorageCosts, bundletypes.StorageCost{StorageProviderId: 2, Cost: math.LegacyMustNewDecFromStr("0.9")})
		s.App().BundlesKeeper.SetParams(s.Ctx(), params)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.CurrentStorageProviderId = storageProviderId
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(200 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert uploader self delegation rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.9)) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(i.ACoins(8829).String()))
		// assert commission rewards
		// (10_000 - (10_000 * 0.01) - (100 * 0.9)) * 0.1 + (100 * 0.9)
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(i.ACoins(1071).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(i.ACoins(100*i.T_KYVE - amountPerBundle).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators, multiple coins and no foreign delegations", func() {
		// ARRANGE
		// fund additionally to the already funded 100acoins
		s.RunTxPoolSuccess(&fundersTypes.MsgFundPool{
			Creator:          i.ALICE,
			Amounts:          sdk.NewCoins(i.BCoin(200*i.T_KYVE), i.CCoin(300*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.BCoin(2*amountPerBundle), i.CCoin(3*amountPerBundle)),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(200 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (coin_amount_per_bundle - (coin_amount_per_bundle * 0.01) - _((100 * 0.5) / (3 * coin_weight))_) * 0.1 + _((100 * 0.5) / (3 * coin_weight))_
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(sdk.NewCoins(i.ACoin(1004), i.BCoin(1987), i.CCoin(2974)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (coin_amount_per_bundle - (coin_amount_per_bundle * 0.01) - _((100 * 0.5) / (3 * coin_weight))_) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.ACoin(8896), i.BCoin(17813), i.CCoin(26726)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(Equal(uint64(200)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(Equal(uint64(300)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE-amountPerBundle), i.BCoin(200*i.T_KYVE-2*amountPerBundle), i.CCoin(300*i.T_KYVE-3*amountPerBundle)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators, multiple coins which are not enough for the storage reward and no foreign delegations", func() {
		// ARRANGE
		// set coin weight of ccoin to a very low value so amount_per_bundle of ccoin can not cover the storage reward
		s.App().FundersKeeper.SetParams(s.Ctx(), fundersTypes.NewParams([]*fundersTypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(2),
			},
			{
				CoinDenom:                 i.C_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyMustNewDecFromStr("0.0000001"),
			},
		}, 0))

		// fund additionally to the already funded 100acoins
		s.RunTxPoolSuccess(&fundersTypes.MsgFundPool{
			Creator:          i.ALICE,
			Amounts:          sdk.NewCoins(i.BCoin(100*i.T_KYVE), i.CCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.BCoin(amountPerBundle), i.CCoin(amountPerBundle)),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(200 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (3 * coin_weight))_) * 0.1 + _((100 * 0.5) / (3 * coin_weight))_
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(sdk.NewCoins(i.ACoin(1004), i.BCoin(997), i.CCoin(9900)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (3 * coin_weight))_) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.ACoin(8896), i.BCoin(8903)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(Equal(uint64(100)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE-amountPerBundle), i.BCoin(100*i.T_KYVE-amountPerBundle), i.CCoin(100*i.T_KYVE-amountPerBundle)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validator, multiple coins and no foreign delegations where one coin is removed from the whitelist", func() {
		// ARRANGE
		// fund additionally to the already funded 100acoins
		s.RunTxPoolSuccess(&fundersTypes.MsgFundPool{
			Creator:          i.ALICE,
			Amounts:          sdk.NewCoins(i.BCoin(100*i.T_KYVE), i.CCoin(100*i.T_KYVE)),
			AmountsPerBundle: sdk.NewCoins(i.BCoin(amountPerBundle), i.CCoin(amountPerBundle)),
		})

		// remove ccoin from whitelist
		s.App().FundersKeeper.SetParams(s.Ctx(), fundersTypes.NewParams([]*fundersTypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.A_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(1),
			},
			{
				CoinDenom:                 i.B_DENOM,
				MinFundingAmount:          math.NewIntFromUint64(10 * i.KYVE),
				MinFundingAmountPerBundle: math.NewInt(amountPerBundle),
				CoinWeight:                math.LegacyNewDec(2),
			},
		}, 0))

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      100,
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(200 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (2 * coin_weight))_) * 0.1 + _((100 * 0.5) / (2 * coin_weight))_
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(sdk.NewCoins(i.ACoin(1012), i.BCoin(1000)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (10_000 - (10_000 * 0.01) - _((100 * 0.5) / (2 * coin_weight))_) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.ACoin(8888), i.BCoin(8900)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM).Uint64()).To(Equal(uint64(100)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM).Uint64()).To(BeZero())

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.ACoin(100*i.T_KYVE-amountPerBundle), i.BCoin(100*i.T_KYVE-amountPerBundle), i.CCoin(100*i.T_KYVE)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})

	It("Produce a valid bundle with multiple validators, multiple coins and real world values", func() {
		// ARRANGE
		// prices are from the 06/06/24

		// Irys: https://node1.bundlr.network/price/1048576 -> 1048 winston/byte * 40 USD/AR * 1.5 / 10**12
		bundleParams := s.App().BundlesKeeper.GetParams(s.Ctx())
		bundleParams.StorageCosts = []bundletypes.StorageCost{
			{StorageProviderId: 1, Cost: math.LegacyMustNewDecFromStr("0.000000006288")},
		}
		s.App().BundlesKeeper.SetParams(s.Ctx(), bundleParams)

		// tkyve -> $KYVE
		// acoin -> $TIA
		// bcoin -> $ARCH
		// ccoin -> $OSMO
		s.App().FundersKeeper.SetParams(s.Ctx(), fundersTypes.NewParams([]*fundersTypes.WhitelistCoinEntry{
			{
				CoinDenom:                 globalTypes.Denom,
				CoinDecimals:              6,
				MinFundingAmount:          math.NewIntFromUint64(100_000_000),    // 100 $KYVE
				MinFundingAmountPerBundle: math.NewInt(100_000),                  // 0.1 $KYVE
				CoinWeight:                math.LegacyMustNewDecFromStr("0.055"), // 0.055 $USD
			},
			{
				CoinDenom:                 i.A_DENOM,
				CoinDecimals:              6,
				MinFundingAmount:          math.NewIntFromUint64(100_000_000),    // 100 $TIA
				MinFundingAmountPerBundle: math.NewInt(100_000),                  // 0.1 $TIA
				CoinWeight:                math.LegacyMustNewDecFromStr("10.33"), // 10.33 $USD
			},
			{
				CoinDenom:                 i.B_DENOM,
				CoinDecimals:              18,
				MinFundingAmount:          s.MustNewIntFromStr("100000000000000000000"), // 100 $ARCH
				MinFundingAmountPerBundle: s.MustNewIntFromStr("100000000000000000"),    // 0.1 $ARCH
				CoinWeight:                math.LegacyMustNewDecFromStr("0.084"),        // 0.084 $USD
			},
			{
				CoinDenom:                 i.C_DENOM,
				CoinDecimals:              6,
				MinFundingAmount:          math.NewIntFromUint64(100_000_000),   // 100 $OSMO
				MinFundingAmountPerBundle: math.NewInt(100_000),                 // 0.1 $OSMO
				CoinWeight:                math.LegacyMustNewDecFromStr("0.84"), // 0.84 $USD
			},
		}, 0))

		// mint another 1,000 $ARCH to Alice
		err := s.MintCoins(i.ALICE, sdk.NewCoins(sdk.NewCoin(i.B_DENOM, s.MustNewIntFromStr("1000000000000000000000"))))
		_ = err

		// defund everything so we have a clean state
		s.RunTxPoolSuccess(&fundersTypes.MsgDefundPool{
			Creator: i.ALICE,
			Amounts: i.ACoins(100 * i.T_KYVE),
		})

		// fund with every coin 100 units for amount and 1 for amount per bundle
		s.RunTxPoolSuccess(&fundersTypes.MsgFundPool{
			Creator:          i.ALICE,
			Amounts:          sdk.NewCoins(i.KYVECoin(100_000_000), i.ACoin(100_000_000), sdk.NewCoin(i.B_DENOM, s.MustNewIntFromStr("100000000000000000000")), i.CCoin(100_000_000)),
			AmountsPerBundle: sdk.NewCoins(i.KYVECoin(1_000_000), i.ACoin(1_000_000), sdk.NewCoin(i.B_DENOM, s.MustNewIntFromStr("1000000000000000000")), i.CCoin(1_000_000)),
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_0_A,
			Staker:        i.STAKER_0,
			PoolId:        0,
			StorageId:     "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			DataSize:      1048576, // 1MB
			DataHash:      "test_hash",
			FromIndex:     0,
			BundleSize:    100,
			FromKey:       "0",
			ToKey:         "99",
			BundleSummary: "test_value",
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.POOL_ADDRESS_1_A,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_VALID,
		})

		initialBalanceStaker1 = s.GetCoinsFromAddress(i.STAKER_1)
		initialBalancePoolAddress1 = s.GetCoinsFromAddress(i.POOL_ADDRESS_1_A)

		c1 := s.GetCoinsFromCommunityPool()

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:       i.POOL_ADDRESS_1_A,
			Staker:        i.STAKER_1,
			PoolId:        0,
			StorageId:     "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			DataSize:      100,
			DataHash:      "test_hash2",
			FromIndex:     100,
			BundleSize:    100,
			FromKey:       "100",
			ToKey:         "199",
			BundleSummary: "test_value2",
		})

		// ASSERT
		// check if bundle got finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal("99"))
		Expect(pool.CurrentSummary).To(Equal("test_value"))
		Expect(pool.CurrentIndex).To(Equal(uint64(100)))
		Expect(pool.TotalBundles).To(Equal(uint64(1)))

		// check if finalized bundle got saved
		finalizedBundle, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeTrue())

		Expect(finalizedBundle.PoolId).To(Equal(uint64(0)))
		Expect(finalizedBundle.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(finalizedBundle.Uploader).To(Equal(i.STAKER_0))
		Expect(finalizedBundle.FromIndex).To(Equal(uint64(0)))
		Expect(finalizedBundle.ToIndex).To(Equal(uint64(100)))
		Expect(finalizedBundle.FromKey).To(Equal("0"))
		Expect(finalizedBundle.ToKey).To(Equal("99"))
		Expect(finalizedBundle.BundleSummary).To(Equal("test_value"))
		Expect(finalizedBundle.DataHash).To(Equal("test_hash"))
		Expect(finalizedBundle.FinalizedAt).NotTo(BeZero())
		Expect(finalizedBundle.StakeSecurity.ValidVotePower).To(Equal(200 * i.KYVE))
		Expect(finalizedBundle.StakeSecurity.TotalVotePower).To(Equal(200 * i.KYVE))

		// check if next bundle proposal got registered
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg"))
		Expect(bundleProposal.Uploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.NextUploader).NotTo(BeEmpty())
		Expect(bundleProposal.DataSize).To(Equal(uint64(100)))
		Expect(bundleProposal.DataHash).To(Equal("test_hash2"))
		Expect(bundleProposal.BundleSize).To(Equal(uint64(100)))
		Expect(bundleProposal.FromKey).To(Equal("100"))
		Expect(bundleProposal.ToKey).To(Equal("199"))
		Expect(bundleProposal.BundleSummary).To(Equal("test_value2"))
		Expect(bundleProposal.UpdatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		poolAccountUploader, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(poolAccountUploader.Points).To(BeZero())

		balanceUploaderValaddress := s.GetCoinsFromAddress(poolAccountUploader.PoolAddress)
		Expect(balanceUploaderValaddress).To(Equal(initialBalancePoolAddress0))

		// check voter status
		poolAccountVoter, _ := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_1, 0)
		Expect(poolAccountVoter.Points).To(BeZero())

		balanceVoterPoolAddress := s.GetCoinsFromAddress(poolAccountVoter.PoolAddress)
		Expect(balanceVoterPoolAddress).To(Equal(initialBalancePoolAddress1))

		balanceVoter := s.GetCoinsFromAddress(poolAccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		// check uploader rewards
		balanceUploader := s.GetCoinsFromAddress(poolAccountUploader.Staker)

		// assert payout transfer
		Expect(balanceUploader.String()).To(Equal(initialBalanceStaker0.String()))
		// assert commission rewards (here we round down since the result of commission rewards gets truncated)
		// (amount_per_bundle - treasury_reward - storage_cost) * uploader_commission + storage_cost
		// storage_cost = 1MB * storage_price / coin_length * coin_price
		// (amount_per_bundle - (amount_per_bundle * 0.01) - _((1048576 * 0.000000006288 * 10**coin_decimals) / (4 * coin_weight))_) * 0.1 + _((1048576 * 0.000000006288) / (4 * coin_weight))_
		Expect(s.App().StakersKeeper.GetOutstandingCommissionRewards(s.Ctx(), poolAccountUploader.Staker).String()).To(Equal(sdk.NewCoins(i.KYVECoin(125_973), i.ACoin(99_143), i.BCoin(116_661_015_771_428_571), i.CCoin(100_765)).String()))
		// assert uploader self delegation rewards (here we round up since the result of delegation rewards is the remainder minus the truncated commission rewards)
		// (amount_per_bundle - (amount_per_bundle * 0.01) - _((29970208 * 0.000000006288 * 1**coin_decimals) / (4 * coin_weight))_) * (1 - 0.1)
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0).String()).To(Equal(sdk.NewCoins(i.KYVECoin(864_027), i.ACoin(890_857), i.BCoin(873_338_984_228_571_429), i.CCoin(889_235)).String()))

		fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

		// assert treasury payout
		c2 := s.GetCoinsFromCommunityPool()
		Expect(c2.Sub(c1...).AmountOf(i.A_DENOM)).To(Equal(math.NewInt(10000)))
		Expect(c2.Sub(c1...).AmountOf(i.B_DENOM)).To(Equal(math.NewInt(10000000000000000)))
		Expect(c2.Sub(c1...).AmountOf(i.C_DENOM)).To(Equal(math.NewInt(10000)))

		// assert total pool funds
		Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId).String()).To(Equal(sdk.NewCoins(i.KYVECoin(99_000_000), i.ACoin(99_000_000), sdk.NewCoin(i.B_DENOM, s.MustNewIntFromStr("99000000000000000000")), i.CCoin(99_000_000)).String()))
		Expect(fundingState.ActiveFunderAddresses).To(HaveLen(1))
	})
})
