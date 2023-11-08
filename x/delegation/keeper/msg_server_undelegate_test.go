package keeper_test

import (
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_undelegate.go

* Undelegate more $KYVE than allowed
* Start undelegation; Check unbonding queue state
* Start undelegation and await unbonding
* Redelegation during undelegation unbonding
* Undelegate Slashed Amount
* Delegate twice and undelegate twice
* Delegate twice and undelegate twice and await unbonding
* Undelegate all after rewards and slashing
* JoinA, Slash, JoinB, PayoutReward
* Slash twice
* Start unbonding, slash twice, payout, await undelegation
*/

var _ = Describe("msg_server_undelegate.go", Ordered, func() {
	s := i.NewCleanChain()

	const aliceSelfDelegation = 100 * i.KYVE
	const bobSelfDelegation = 100 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		CreateFundedPool(s)

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  aliceSelfDelegation,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  bobSelfDelegation,
		})

		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 10_000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.BOB,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
			Amount:     0,
		})

		_, aliceFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.ALICE)
		Expect(aliceFound).To(BeTrue())

		_, bobFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.BOB)
		Expect(bobFound).To(BeTrue())

		s.CommitAfterSeconds(7)
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(s)
	})

	It("Undelegate more $KYVE than allowed", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))

		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorError(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  11 * i.KYVE,
		})

		// ASSERT
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(s.Ctx(), i.DUMMY[0])).To(BeEmpty())
	})

	It("Start undelegation; Check unbonding queue state", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		s.CommitAfterSeconds(1)

		// ASSERT

		// Delegation amount stays the same (due to unbonding)
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		unbondingEntries := s.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(s.Ctx(), i.DUMMY[0])
		Expect(unbondingEntries).To(HaveLen(1))
		Expect(unbondingEntries[0].Staker).To(Equal(i.ALICE))
		Expect(unbondingEntries[0].Delegator).To(Equal(i.DUMMY[0]))
		Expect(unbondingEntries[0].Amount).To(Equal(5 * i.KYVE))
		Expect(unbondingEntries[0].CreationTime).To(Equal(uint64(s.Ctx().BlockTime().Unix() - 1)))
	})

	It("Start undelegation and await unbonding", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()) + 1)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(995 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 5*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(5 * i.KYVE))

		unbondingEntries := s.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(s.Ctx(), i.DUMMY[0])
		Expect(unbondingEntries).To(BeEmpty())
	})

	It("Redelegation during undelegation unbonding", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.BOB,
			Amount:     10 * i.KYVE,
		})

		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()) + 1)
		s.CommitAfterSeconds(1)

		// ASSERT

		// Unbonding should have had no effect
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(0 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.BOB, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		unbondingEntries := s.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntriesOfDelegator(s.Ctx(), i.DUMMY[0])
		Expect(unbondingEntries).To(BeEmpty())
	})

	It("Undelegate Slashed Amount", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		params := s.App().DelegationKeeper.GetParams(s.Ctx())
		params.UploadSlash = sdk.MustNewDecFromStr("0.1")
		s.App().DelegationKeeper.SetParams(s.Ctx(), params)
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)

		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(9 * i.KYVE))

		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()) + 1)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(999 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(0 * i.KYVE))
	})

	It("Delegate twice and undelegate twice", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.PerformValidityChecks()

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(980 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 30*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(20 * i.KYVE))

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  8 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		// ASSERT
		unbondingEntries := s.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntries(s.Ctx())

		Expect(unbondingEntries).To(HaveLen(2))
		Expect(unbondingEntries[0].Staker).To(Equal(i.ALICE))
		Expect(unbondingEntries[0].Delegator).To(Equal(i.DUMMY[0]))
		Expect(unbondingEntries[0].Amount).To(Equal(5 * i.KYVE))
		Expect(unbondingEntries[0].CreationTime).To(Equal(uint64(s.Ctx().BlockTime().Unix() - 20)))

		Expect(unbondingEntries[1].Staker).To(Equal(i.ALICE))
		Expect(unbondingEntries[1].Delegator).To(Equal(i.DUMMY[1]))
		Expect(unbondingEntries[1].Amount).To(Equal(8 * i.KYVE))
		Expect(unbondingEntries[1].CreationTime).To(Equal(uint64(s.Ctx().BlockTime().Unix() - 10)))
	})

	It("Delegate twice and undelegate twice and await unbonding", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(980 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 30*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(20 * i.KYVE))

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  8 * i.KYVE,
		})

		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()) + 1)
		s.CommitAfterSeconds(1)

		// ASSERT
		unbondingEntries := s.App().DelegationKeeper.GetAllUnbondingDelegationQueueEntries(s.Ctx())
		Expect(unbondingEntries).To(BeEmpty())

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(995 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(5 * i.KYVE))

		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(988 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(12 * i.KYVE))
	})

	It("Undelegate all after rewards and slashing", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(10 * i.KYVE))

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(980 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 30*i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(20 * i.KYVE))

		s.PerformValidityChecks()

		// Payout rewards
		// Alice: 100   100/130 * 10 * 1e9 = 7_692_307_692
		// Dummy0: 10   10/130 * 10 * 1e9 = 769_230_769
		// Dummy1: 20   20/130 * 10 * 1e9 = 1_538_461_538
		PayoutRewards(s, i.ALICE, 10*i.KYVE)

		// Collect
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
		})

		// Slash 10%
		params := s.App().DelegationKeeper.GetParams(s.Ctx())
		params.UploadSlash = sdk.MustNewDecFromStr("0.1")
		s.App().DelegationKeeper.SetParams(s.Ctx(), params)
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)

		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(9 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(18 * i.KYVE))

		// ACT
		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  9 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  18 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()) + 1)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(999*i.KYVE + uint64(769_230_769)))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(998*i.KYVE + uint64(1_538_461_538)))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(uint64(float64(aliceSelfDelegation) * 0.9)))

		delegationEntries := s.App().DelegationKeeper.GetAllDelegationEntries(s.Ctx())
		delegators := s.App().DelegationKeeper.GetAllDelegators(s.Ctx())
		slashes := s.App().DelegationKeeper.GetAllDelegationSlashEntries(s.Ctx())

		Expect(len(slashes)).To(Equal(1))
		Expect(len(delegationEntries)).To(Equal(4))
		Expect(delegators).To(HaveLen(2))
	})

	It("JoinA, Slash, JoinB, PayoutReward", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		params := s.App().DelegationKeeper.GetParams(s.Ctx())
		params.UploadSlash = sdk.MustNewDecFromStr("0.5")
		s.App().DelegationKeeper.SetParams(s.Ctx(), params)
		s.PerformValidityChecks()

		// Slash 50%
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		// Dummy0: 5$KYVE Dummy1: 20$KYVE
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal((50 + 25) * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(5 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(20 * i.KYVE))

		// ACT

		// Alice: 50    50 / 75 * 10 * 1e9 = 6_666_666_666
		// Dummy0: 5    5 / 75 * 10 * 1e9 = 666_666_666
		// Dummy1: 20   20 / 75 * 10 * 1e9 = 2_666_666_666
		PayoutRewards(s, i.ALICE, 10*i.KYVE)

		// ASSERT
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(666_666_666)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(2_666_666_666)))

		// must be the same as before
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal((50 + 25) * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(5 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(20 * i.KYVE))
	})

	It("Slash twice", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.PerformValidityChecks()

		// ACT
		params := s.App().DelegationKeeper.GetParams(s.Ctx())
		params.UploadSlash = sdk.MustNewDecFromStr("0.5")
		s.App().DelegationKeeper.SetParams(s.Ctx(), params)

		// Slash 50% twice
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)

		// ASSERT
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(25*i.KYVE + uint64(2_500_000_000+5_000_000_000)))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(2_500_000_000)))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(5_000_000_000)))
	})

	It("Start unbonding, slash twice, payout, await undelegation", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.CommitAfterSeconds(10)

		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgUndelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.PerformValidityChecks()

		// ACT
		params := s.App().DelegationKeeper.GetParams(s.Ctx())
		params.UploadSlash = sdk.MustNewDecFromStr("0.5")
		s.App().DelegationKeeper.SetParams(s.Ctx(), params)
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)

		// Alice: 25    25 / 32.5 * 1e10 = 7_692_307_692
		// Dummy0: 2.5  2.5 / 32.5 * 1e10 = 769_230_769
		// Dummy1: 5    5 / 32.5 * 1e10 = 1_538_461_538
		PayoutRewards(s, i.ALICE, 10*i.KYVE)

		s.CommitAfterSeconds(s.App().DelegationKeeper.GetUnbondingDelegationTime(s.Ctx()) + 1)
		s.CommitAfterSeconds(1)

		// ASSERT
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(25 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeZero())
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.ALICE, i.DUMMY[1])).To(BeZero())

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(uint64(1000e9 - 7_500_000_000 + 769_230_769)))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(uint64(1000e9 - 15_000_000_000 + 1_538_461_538)))
	})
})
