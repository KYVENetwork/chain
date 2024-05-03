package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_withdraw_rewards.go

* Payout rewards which cause rounding issues and withdraw
* Withdraw from a non-existing delegator
* Test invalid payouts to delegators
* Withdraw rewards which are zero
* Withdraw rewards with multiple slashes
* Payout rewards with multiple denoms
* Payout mixed rewards multiple times and multiple delegation steps
* Withdraw multiple denoms, one after the other
*/

var _ = Describe("msg_server_withdraw_rewards.go", Ordered, func() {
	s := i.NewCleanChain()

	const aliceSelfDelegation = 0 * i.KYVE
	const bobSelfDelegation = 0 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		CreateFundedPool(s)

		// Stake
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  aliceSelfDelegation,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  bobSelfDelegation,
		})

		_, stakerFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.ALICE)
		Expect(stakerFound).To(BeTrue())

		_, stakerFound = s.App().StakersKeeper.GetStaker(s.Ctx(), i.BOB)
		Expect(stakerFound).To(BeTrue())
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(s)
	})

	It("Payout rewards which cause rounding issues and withdraw", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(990 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(990 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 30*i.KYVE))

		delegationModuleBalanceBefore := s.GetBalanceFromModule(types.ModuleName)
		fundersModuleBalanceBefore := s.GetBalanceFromModule(funderstypes.ModuleName)
		s.PerformValidityChecks()

		// ACT

		// Alice: 100
		// Dummy0: 10
		// Dummy1: 0
		PayoutRewards(s, i.ALICE, sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(20*i.KYVE))))

		// ASSERT
		delegationModuleBalanceAfter := s.GetBalanceFromModule(types.ModuleName)
		fundersModuleBalanceAfter := s.GetBalanceFromModule(funderstypes.ModuleName)

		Expect(delegationModuleBalanceAfter).To(Equal(delegationModuleBalanceBefore + 20*i.KYVE))
		Expect(fundersModuleBalanceAfter).To(Equal(fundersModuleBalanceBefore - 20*i.KYVE))

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(6666666666)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(6666666666)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[2]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(uint64(6666666666)))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(uint64(996666666666)))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(uint64(996666666666)))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(uint64(996666666666)))

		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal(uint64(30000000002)))
	})

	It("Withdraw from a non-existing delegator", func() {
		// ARRANGE
		balanceDummy1Before := s.GetBalanceFromAddress(i.DUMMY[0])
		balanceCharlieBefore := s.GetBalanceFromAddress(i.CHARLIE)
		balanceAliceBefore := s.GetBalanceFromAddress(i.ALICE)
		delegationBalance := s.GetBalanceFromModule(types.ModuleName)

		s.PerformValidityChecks()

		// ACT
		s.RunTxDelegatorError(&types.MsgWithdrawRewards{
			Creator: i.CHARLIE,
			Staker:  i.ALICE,
		})

		s.RunTxDelegatorError(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.CHARLIE,
		})

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(balanceDummy1Before))
		Expect(s.GetBalanceFromAddress(i.CHARLIE)).To(Equal(balanceCharlieBefore))
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(balanceAliceBefore))
		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal(delegationBalance))
	})

	It("Test invalid payouts to delegators", func() {
		// ARRANGE

		// fund pool module
		s.RunTxFundersSuccess(&funderstypes.MsgFundPool{
			Creator:         i.ALICE,
			PoolId:          0,
			Amount:          100 * i.KYVE,
			AmountPerBundle: 1 * i.KYVE,
		})

		// ACT
		// not enough balance in pool module
		err1 := s.App().DelegationKeeper.PayoutRewards(s.Ctx(), i.ALICE, sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(20000*i.KYVE))), pooltypes.ModuleName)
		// staker does not exist
		err2 := s.App().DelegationKeeper.PayoutRewards(s.Ctx(), i.DUMMY[20], sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(1*i.KYVE))), pooltypes.ModuleName)

		// ASSERT
		Expect(err1).To(HaveOccurred())
		Expect(err2).To(HaveOccurred())
	})

	It("Withdraw rewards which are zero", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  1,
		})
		startBalance := s.GetBalanceFromAddress(i.DUMMY[0])

		// ACT
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeEmpty())

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		// ASSERT
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(startBalance))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeEmpty())
	})

	It("Withdraw rewards with multiple slashes", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		startBalance := s.GetBalanceFromAddress(i.DUMMY[0])

		// ACT
		params := s.App().DelegationKeeper.GetParams(s.Ctx())
		params.UploadSlash = math.LegacyMustNewDecFromStr("0.5")
		s.App().DelegationKeeper.SetParams(s.Ctx(), params)

		// Slash 50%
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)
		PayoutRewards(s, i.ALICE, sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(5*i.KYVE))))

		// Slash 50% again
		s.App().DelegationKeeper.SlashDelegators(s.Ctx(), 0, i.ALICE, types.SLASH_TYPE_UPLOAD)
		PayoutRewards(s, i.ALICE, sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(5*i.KYVE))))

		s.PerformValidityChecks()

		// ASSERT
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).AmountOf(globalTypes.Denom).Uint64()).To(Equal(10 * i.KYVE))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeEmpty())
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(startBalance + 10*i.KYVE))
	})

	It("Payout rewards with multiple denoms", func() {
		mintErr1 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "acoin")
		Expect(mintErr1).To(BeNil())

		mintErr2 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "bcoin")
		Expect(mintErr2).To(BeNil())

		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(990 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(990 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 30*i.KYVE))

		delegationModuleBalanceBefore := s.GetCoinsFromModule(types.ModuleName)
		poolsModuleBalanceBefore := s.GetCoinsFromModule(pooltypes.ModuleName)
		s.PerformValidityChecks()

		// ACT
		payoutCoins := sdk.NewCoins(
			sdk.NewCoin("acoin", math.NewInt(10*1_000_000)),
			sdk.NewCoin("bcoin", math.NewInt(5*1_000_000)),
		)
		PayoutRewards(s, i.ALICE, payoutCoins)

		// ASSERT
		delegationModuleBalanceAfter := s.GetCoinsFromModule(types.ModuleName)
		poolsModuleBalanceAfter := s.GetCoinsFromModule(pooltypes.ModuleName)

		Expect(delegationModuleBalanceAfter).To(Equal(delegationModuleBalanceBefore.Add(payoutCoins...)))
		Expect(poolsModuleBalanceAfter).To(Equal(poolsModuleBalanceBefore.Sub(payoutCoins...)))

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).String()).To(Equal("3333333acoin,1666666bcoin"))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).String()).To(Equal("3333333acoin,1666666bcoin"))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[2]).String()).To(Equal("3333333acoin,1666666bcoin"))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
		})

		Expect(s.GetCoinsFromAddress(i.DUMMY[0]).String()).To(Equal("3333333acoin,1666666bcoin,990000000000tkyve"))
		Expect(s.GetCoinsFromAddress(i.DUMMY[1]).String()).To(Equal("3333333acoin,1666666bcoin,990000000000tkyve"))
		Expect(s.GetCoinsFromAddress(i.DUMMY[2]).String()).To(Equal("3333333acoin,1666666bcoin,990000000000tkyve"))
		Expect(s.GetCoinsFromModule(types.ModuleName).String()).To(Equal("1acoin,2bcoin,30000000000tkyve"))
	})

	It("Payout mixed rewards multiple times and multiple delegation steps", func() {
		mintErr1 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "acoin")
		Expect(mintErr1).To(BeNil())

		mintErr2 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "bcoin")
		Expect(mintErr2).To(BeNil())

		mintErr3 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "ccoin")
		Expect(mintErr3).To(BeNil())

		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(980 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(995 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(1000 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 25*i.KYVE))

		delegationModuleBalanceBefore := s.GetCoinsFromModule(types.ModuleName)
		poolsModuleBalanceBefore := s.GetCoinsFromModule(pooltypes.ModuleName)
		s.PerformValidityChecks()

		// ACT
		payoutCoins1 := sdk.NewCoins(
			sdk.NewCoin("acoin", math.NewInt(10*1_000_000)),
			sdk.NewCoin("bcoin", math.NewInt(5*1_000_000)),
		)
		PayoutRewards(s, i.ALICE, payoutCoins1)

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
			Amount:  3 * i.KYVE,
		})

		payoutCoins2 := sdk.NewCoins(
			sdk.NewCoin("bcoin", math.NewInt(8*1_000_000)),
			sdk.NewCoin("ccoin", math.NewInt(7*1_000_000)),
		)
		PayoutRewards(s, i.ALICE, payoutCoins2)

		// ASSERT
		delegationModuleBalanceAfter := s.GetCoinsFromModule(types.ModuleName)
		poolsModuleBalanceAfter := s.GetCoinsFromModule(pooltypes.ModuleName)

		Expect(delegationModuleBalanceAfter).To(Equal(
			delegationModuleBalanceBefore.Add(payoutCoins1.Add(payoutCoins2...).Add(sdk.NewInt64Coin(globalTypes.Denom, int64(3*i.KYVE)))...)),
		)
		Expect(poolsModuleBalanceAfter).To(Equal(poolsModuleBalanceBefore.Sub(payoutCoins1.Add(payoutCoins2...)...)))

		// Calculate outstanding rewards
		// phase1:
		//   reward: 10_000_000 acoin, 5_000_000 bcoin
		//   delegation_shares: d0: 20, d1: 5
		//   => d0: 20/25 * 10_000_000 = 8_000_000 acoin; 20/25 * 5_000_000 = 4_000_000 bcoin
		//   => d1: 5/25 * 10_000_000 = 2_000_000 acoin; 5/25 * 5_000_000 = 1_000_000 bcoin
		// phase2:
		//   reward: 8_000_000 bcoin, 7_000_000 ccoin
		//   delegation_shares: d0: 20, d1: 5 d2: 3
		//   => d0: 20/28 * 8_000_000 = 5_714_285 bcoin; 20/28 * 7_000_000 = 5_000_000 ccoin
		//   => d1: 5/28 * 8_000_000 = 1_428_571 bcoin; 5/28 * 7_000_000 = 1_250_000 ccoin
		//   => d2: 3/28 * 8_000_000 = 857_142 bcoin; 3/28 * 7_000_000 = 750_000 ccoin
		// SUM: d0: 8_000_000 acoin, 9_714_285 bcoin, 5_000_000 ccoin,
		// SUM: d1: 2_000_000 acoin, 2_428_571 bcoin, 1_250_000 ccoin
		// SUM: d2: 0 acoin, 857_142 bcoin, 750_000 ccoin
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).String()).To(Equal("8000000acoin,9714285bcoin,5000000ccoin"))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).String()).To(Equal("2000000acoin,2428571bcoin,1250000ccoin"))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[2]).String()).To(Equal("857142bcoin,750000ccoin"))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
		})

		Expect(s.GetCoinsFromAddress(i.DUMMY[0]).String()).To(Equal("8000000acoin,9714285bcoin,5000000ccoin,980000000000tkyve"))
		Expect(s.GetCoinsFromAddress(i.DUMMY[1]).String()).To(Equal("2000000acoin,2428571bcoin,1250000ccoin,995000000000tkyve"))
		Expect(s.GetCoinsFromAddress(i.DUMMY[2]).String()).To(Equal("857142bcoin,750000ccoin,997000000000tkyve"))
		Expect(s.GetCoinsFromModule(types.ModuleName).String()).To(Equal("2bcoin,28000000000tkyve"))
	})

	It("Withdraw multiple denoms, one after the other", func() {
		mintErr1 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "acoin")
		Expect(mintErr1).To(BeNil())

		mintErr2 := s.MintDenomToModule(pooltypes.ModuleName, 1000*1_000_000, "bcoin")
		Expect(mintErr2).To(BeNil())

		// ARRANGE
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  20 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  5 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(980 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(995 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(1000 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 25*i.KYVE))

		delegationModuleBalanceBefore := s.GetCoinsFromModule(types.ModuleName)
		poolsModuleBalanceBefore := s.GetCoinsFromModule(pooltypes.ModuleName)
		s.PerformValidityChecks()

		// ACT
		payoutCoins1 := sdk.NewCoins(
			sdk.NewCoin("acoin", math.NewInt(10*1_000_000)),
		)
		PayoutRewards(s, i.ALICE, payoutCoins1)

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		payoutCoins2 := sdk.NewCoins(
			sdk.NewCoin("bcoin", math.NewInt(8*1_000_000)),
		)
		PayoutRewards(s, i.ALICE, payoutCoins2)

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
		})

		// ASSERT
		delegationModuleBalanceAfter := s.GetCoinsFromModule(types.ModuleName)
		poolsModuleBalanceAfter := s.GetCoinsFromModule(pooltypes.ModuleName)

		compare := delegationModuleBalanceBefore.Add(payoutCoins1...)
		compare = compare.Add(payoutCoins2...)
		compare = compare.Add(sdk.NewInt64Coin(globalTypes.Denom, int64(10*i.KYVE)))
		compare = compare.Sub(sdk.NewInt64Coin("acoin", 8_000_000))
		compare = compare.Sub(sdk.NewInt64Coin("bcoin", 2_285_714))
		Expect(delegationModuleBalanceAfter).To(Equal(compare))
		Expect(poolsModuleBalanceAfter).To(Equal(poolsModuleBalanceBefore.Sub(payoutCoins1.Add(payoutCoins2...)...)))

		// Calculate outstanding rewards
		// phase1:
		//   reward: 10_000_000 acoin
		//   delegation_shares: d0: 20, d1: 5
		//   => d0: 20/25 * 10_000_000 = 8_000_000 acoin
		//   => d1: 5/25 * 10_000_000 = 2_000_000 acoin
		// d0: withdraw rewards (-8_000_000 acoin)
		// phase2:
		//   reward: 8_000_000 bcoin
		//   delegation_shares: d0: 20, d1: 5 d2: 10
		//   => d0: 20/35 * 8_000_000 = 4_571_428 bcoin
		//   => d1: 5/35 * 8_000_000 = 1_142_857 bcoin
		//   => d2: 10/35 * 8_000_000 = 2_285_714 bcoin
		// d2: withdraw rewards (-2_285_714 bcoin)
		// SUM: d0: 4_571_428 bcoin
		// SUM: d1: 2_000_000 acoin, 1_142_857 bcoin
		// SUM: d2: 0
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0]).String()).To(Equal("4571428bcoin"))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1]).String()).To(Equal("2000000acoin,1142857bcoin"))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[2]).String()).To(Equal(""))

		Expect(s.GetCoinsFromAddress(i.DUMMY[0]).String()).To(Equal("8000000acoin,980000000000tkyve"))
		Expect(s.GetCoinsFromAddress(i.DUMMY[1]).String()).To(Equal("995000000000tkyve"))
		Expect(s.GetCoinsFromAddress(i.DUMMY[2]).String()).To(Equal("2285714bcoin,990000000000tkyve"))
		Expect(s.GetCoinsFromModule(types.ModuleName).String()).To(Equal("2000000acoin,5714286bcoin,35000000000tkyve"))
	})
})
