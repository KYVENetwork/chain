package keeper_test

import (
	"cosmossdk.io/math"
	multicoinrewardstypes "github.com/KYVENetwork/chain/x/multi_coin_rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_leave_pool.go

* MultiCoin disabled, Only native rewards
* MultiCoin disabled, Only foreign rewards
* MultiCoin disabled, Mixed rewards
* MultiCoin enabled, Only native rewards
* MultiCoin enabled, Only foreign rewards
* MultiCoin enabled, Mixed rewards
* MultiCoin empty claim
* MultiCoin claim pending rewards
* MultiCoin enabled, claim, disable, claim again
* Claim after period is over
* Claim rewards indirectly by delegation

*/

var _ = Describe("msg_server_toggle_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var gov string
	var validator1 i.TestValidatorAddress

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
		gov = s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

		// create pool
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		// create staker
		validator1 = s.CreateNewValidator("MyValidator-1", 1000*i.KYVE)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("MultiCoin disabled, Only native rewards", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(100*i.T_KYVE))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100000000tkyve"))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9100000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin disabled, Only foreign rewards", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100acoin,50bcoin"))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal("100acoin,50bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin disabled, Mixed rewards", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100acoin,50bcoin,200000000tkyve"))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal("100acoin,50bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin enabled, Only native rewards", func() {
		// Arrange
		_ = s.App().MultiCoinRewardsKeeper.MultiCoinRewardsEnabled.Set(s.Ctx(), validator1.AccAddress)
		payoutRewards(s, validator1.Address, i.KYVECoins(100*i.T_KYVE))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100000000tkyve"))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9100000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin enabled, Only foreign rewards", func() {
		// Arrange
		_ = s.App().MultiCoinRewardsKeeper.MultiCoinRewardsEnabled.Set(s.Ctx(), validator1.AccAddress)
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100acoin,50bcoin"))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000100acoin,10000000050bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin enabled, Mixed rewards", func() {
		// Arrange
		_ = s.App().MultiCoinRewardsKeeper.MultiCoinRewardsEnabled.Set(s.Ctx(), validator1.AccAddress)
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName)).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100acoin,50bcoin,200000000tkyve"))

		// ACT
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000100acoin,10000000050bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin empty claim", func() {
		// Arrange
		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal(""))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		// ACT
		s.RunTxSuccess(&multicoinrewardstypes.MsgToggleMultiCoinRewards{
			Creator: validator1.Address,
			Enabled: true,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin claim pending rewards", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal("100acoin,50bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		// ACT
		s.RunTxSuccess(&multicoinrewardstypes.MsgToggleMultiCoinRewards{
			Creator: validator1.Address,
			Enabled: true,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000100acoin,10000000050bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("MultiCoin enabled, claim, disable, claim again", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		s.RunTxSuccess(&multicoinrewardstypes.MsgToggleMultiCoinRewards{
			Creator: validator1.Address,
			Enabled: true,
		})

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000100acoin,10000000050bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		// ACT
		s.RunTxSuccess(&multicoinrewardstypes.MsgToggleMultiCoinRewards{
			Creator: validator1.Address,
			Enabled: false,
		})

		payoutRewards(s, validator1.Address, i.KYVECoins(300*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(150))
		payoutRewards(s, validator1.Address, i.BCoins(70))

		s.CommitAfterSeconds(1)

		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000100acoin,10000000050bcoin,10000000000ccoin,9500000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal("150acoin,70bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("Claim after period is over", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		s.RunTxSuccess(&distributionTypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
		})

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal("100acoin,50bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())

		// ACT
		s.CommitAfterSeconds(s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx()).MultiCoinDistributionPendingTime)
		s.CommitAfterSeconds(1)
		s.RunTxSuccess(&multicoinrewardstypes.MsgToggleMultiCoinRewards{
			Creator: validator1.Address,
			Enabled: true,
		})

		// ASSERT
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9200000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName).String()).To(Equal("100acoin,50bcoin"))
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})

	It("Claim rewards indirectly through another delegation", func() {
		// Arrange
		payoutRewards(s, validator1.Address, i.KYVECoins(200*i.T_KYVE))
		payoutRewards(s, validator1.Address, i.ACoins(100))
		payoutRewards(s, validator1.Address, i.BCoins(50))

		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,9000000000tkyve"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address).String()).To(Equal("100acoin,50bcoin,200000000tkyve"))

		// ACT
		s.CommitAfterSeconds(s.App().MultiCoinRewardsKeeper.GetParams(s.Ctx()).MultiCoinDistributionPendingTime)
		s.CommitAfterSeconds(1)
		s.RunTxSuccess(&stakingTypes.MsgDelegate{
			DelegatorAddress: validator1.Address,
			ValidatorAddress: validator1.ValAddress,
			Amount:           i.KYVECoin(400 * i.T_KYVE),
		})

		// ASSERT
		// 9,000 + 200 (rewards) - 400 (delegation) = 8,800
		Expect(s.App().BankKeeper.GetAllBalances(s.Ctx(), validator1.AccAddress).String()).To(Equal("10000000000acoin,10000000000bcoin,10000000000ccoin,8800000000tkyve"))
		cosmosValidator, _ := s.App().StakingKeeper.GetValidator(s.Ctx(), sdk.ValAddress(validator1.AccAddress))
		Expect(cosmosValidator.Tokens.Int64()).To(Equal(1400 * i.T_KYVE))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.ModuleName).String()).To(Equal("100acoin,50bcoin"))
		Expect(s.GetCoinsFromModule(multicoinrewardstypes.MultiCoinRewardsRedistributionAccountName).String()).To(BeEmpty())
		Expect(s.App().StakersKeeper.GetOutstandingRewards(s.Ctx(), validator1.Address, validator1.Address)).To(BeEmpty())
	})
})
