package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	"github.com/KYVENetwork/chain/util"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_join_pool.go

* Test if a newly created staker is participating in no pools yet
* Join the first pool as the first staker to a newly created pool
* Join disabled pool
* Join a pool where other stakers have already joined
* Self-Delegate more KYVE after joining a pool
* Join a pool with the same pool address as the staker address
* Try to join the same pool with the same pool address again
* Try to join the same pool with a different pool address
* Try to join another pool with the same pool address again
* Try to join another pool with a pool address that is already used by another staker
* Try to join another pool with a different pool address
* Try to join pool with empty pool address
* Try to join pool with empty commission
* Try to join pool with empty stake fraction
* Join a pool with a pool address which does not exist on chain yet
* Join a pool with a pool address which does not exist on chain yet and send 0 funds
* Join a pool with an invalid pool address
* Join a pool and fund the pool address with more KYVE than available in balance
* Kick out lowest staker by joining a full pool
* Fail to kick out lowest staker because not enough stake
* Kick out lowest staker with respect to stake + delegation
* Fail to kick out lowest staker because not enough stake + delegation
* Join pool again with same pool address after staker has left pool
* Join pool again with different pool address after staker has left pool

*/

var _ = Describe("msg_server_join_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := uint64(0)
	initialBalancePoolAddress0 := uint64(0)
	initialBalancePoolAddress1 := uint64(0)

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create pool
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		s.SetMaxVotingPower("1")

		s.CreateValidator(i.STAKER_0, "Staker-0", int64(100*i.KYVE))

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalancePoolAddress0 = s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)
		initialBalancePoolAddress1 = s.GetBalanceFromAddress(i.POOL_ADDRESS_0_B)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Test if a newly created staker is participating in no pools yet", func() {
		// ASSERT
		poolAccounts := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccounts).To(HaveLen(0))
	})

	It("Join the first pool as the first staker to a newly created pool", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterPoolAddress0 := s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterPoolAddress0 - initialBalancePoolAddress0).To(Equal(100 * i.KYVE))

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal(i.POOL_ADDRESS_0_A))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join disabled pool", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		pool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		_, err := s.RunTx(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        1,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		Expect(err.Error()).To(Equal("can not join disabled pool: internal logic error"))

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterPoolAddress0 := s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(0 * i.KYVE))
		Expect(balanceAfterPoolAddress0 - initialBalancePoolAddress0).To(Equal(0 * i.KYVE))

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(0))

		_, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 1)

		Expect(active).To(BeFalse())

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 1)

		Expect(poolAccountsOfPool).To(HaveLen(0))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 1)

		Expect(totalStakeOfPool).To(Equal(0 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 1)).To(Equal(0 * i.KYVE))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("join a pool where other stakers have already joined", func() {
		// ARRANGE
		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        0 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        0 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterPoolAddress0 := s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(BeZero())
		Expect(balanceAfterPoolAddress0 - initialBalancePoolAddress0).To(BeZero())

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal(i.POOL_ADDRESS_0_A))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(2))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(200 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Self-Delegate more KYVE after joining a pool", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		// ACT
		s.RunTxSuccess(stakingtypes.NewMsgDelegate(
			i.STAKER_0,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(50*i.KYVE)),
		))

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal(i.POOL_ADDRESS_0_A))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(1))

		totalStakeOfPool = s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(150 * i.KYVE))

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Try to join the same pool with the same pool address again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))
	})

	It("join a pool with the same pool address as the staker address", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.STAKER_0,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Try to join the same pool with a different pool address", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join another pool with the same pool address again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        1,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join pool with a pool address that is already used by another staker", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        1,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join pool with a pool address that is already used by another staker", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(100*i.KYVE))

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_1)
		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Try to join another pool with a different pool address", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        1,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccountsOfStaker).To(HaveLen(2))
	})

	It("Try to join pool with empty pool address", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   "",
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Try to join pool with empty commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Try to join pool with empty stake fraction", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:     i.STAKER_0,
			PoolId:      0,
			PoolAddress: i.POOL_ADDRESS_0_A,
			Amount:      100 * i.KYVE,
			Commission:  math.LegacyMustNewDecFromStr("0.1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Join a pool with a pool address which does not exist on chain yet", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   "kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x",
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterUnknown := s.GetBalanceFromAddress("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x")

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterUnknown).To(Equal(100 * i.KYVE))

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x"))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with a pool address which does not exist on chain yet and send 0 funds", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   "kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x",
			Amount:        0 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterUnknown := s.GetBalanceFromAddress("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x")

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(BeZero())
		Expect(balanceAfterUnknown).To(BeZero())

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x"))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))
		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with an invalid pool address", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   "invalid_valaddress",
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Join a pool and fund the pool address with more KYVE than available in balance", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   "invalid_valaddress",
			Amount:        initialBalanceStaker0 + 1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.ALICE)

		Expect(poolAccountsOfStaker).To(BeEmpty())
	})

	It("Kick out lowest staker by joining a full pool", func() {
		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		for k := 0; k < 49; k++ {
			s.CreateValidator(i.DUMMY[k], fmt.Sprintf("dummy-%d", k), int64(150*i.KYVE))
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:       i.DUMMY[k],
				PoolId:        0,
				PoolAddress:   i.VALDUMMY[k],
				Amount:        1,
				Commission:    math.LegacyMustNewDecFromStr("0.1"),
				StakeFraction: math.LegacyMustNewDecFromStr("1"),
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(150*i.KYVE))

		// Act
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// Assert
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 150) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_0))
	})

	It("Fail to kick out lowest staker because not enough stake", func() {
		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		for k := 0; k < 49; k++ {
			s.CreateValidator(i.DUMMY[k], fmt.Sprintf("dummy-%d", k), int64(150*i.KYVE))
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:       i.DUMMY[k],
				PoolId:        0,
				PoolAddress:   i.VALDUMMY[k],
				Amount:        1,
				Commission:    math.LegacyMustNewDecFromStr("0.1"),
				StakeFraction: math.LegacyMustNewDecFromStr("1"),
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(50*i.KYVE))

		// Act
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// Assert
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_1))
	})

	It("Kick out lowest staker with respect to stake + delegation", func() {
		// ARRANGE
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        1 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		for k := 0; k < 49; k++ {
			s.CreateValidator(i.DUMMY[k], fmt.Sprintf("dummy-%d", k), int64(150*i.KYVE))
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:       i.DUMMY[k],
				PoolId:        0,
				PoolAddress:   i.VALDUMMY[k],
				Amount:        1 * i.KYVE,
				Commission:    math.LegacyMustNewDecFromStr("0.1"),
				StakeFraction: math.LegacyMustNewDecFromStr("1"),
			})
		}

		// Alice is lowest staker and all stakers are full now.
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(150*i.KYVE))

		s.RunTxSuccess(stakingtypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_0),
			sdk.NewInt64Coin(globalTypes.Denom, int64(150*i.KYVE)),
		))

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 250) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).NotTo(ContainElement(i.STAKER_1))
	})

	It("Fail to kick out lowest staker because not enough stake", func() {
		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		for k := 0; k < 49; k++ {
			s.CreateValidator(i.DUMMY[k], fmt.Sprintf("dummy-%d", k), int64(150*i.KYVE))
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:       i.DUMMY[k],
				PoolId:        0,
				PoolAddress:   i.VALDUMMY[k],
				Amount:        1,
				Commission:    math.LegacyMustNewDecFromStr("0.1"),
				StakeFraction: math.LegacyMustNewDecFromStr("1"),
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(50*i.KYVE))

		// Act
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        1,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// Assert
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_1))
	})

	It("Fail to kick out lowest staker because not enough stake + delegation", func() {
		// ARRANGE
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        1 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		for k := 0; k < 49; k++ {
			s.CreateValidator(i.DUMMY[k], fmt.Sprintf("dummy-%d", k), int64(150*i.KYVE))
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:       i.DUMMY[k],
				PoolId:        0,
				PoolAddress:   i.VALDUMMY[k],
				Amount:        1 * i.KYVE,
				Commission:    math.LegacyMustNewDecFromStr("0.1"),
				StakeFraction: math.LegacyMustNewDecFromStr("1"),
			})
		}

		// Alice is lowest staker and all stakers are full now.
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.CreateValidator(i.STAKER_1, "Staker-1", int64(50*i.KYVE))

		s.RunTxSuccess(stakingtypes.NewMsgDelegate(
			i.ALICE,
			util.MustValaddressFromOperatorAddress(i.STAKER_1),
			sdk.NewInt64Coin(globalTypes.Denom, int64(50*i.KYVE)),
		))

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_1,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_1_A,
			Amount:        0,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		Expect(s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).NotTo(ContainElement(i.STAKER_1))
	})

	It("Join pool again with same pool address after staker has left pool", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		_, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(active).To(BeFalse())

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        50 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.2"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterPoolAddress0 := s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(150 * i.KYVE))
		Expect(balanceAfterPoolAddress0 - initialBalancePoolAddress0).To(Equal(150 * i.KYVE))

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal(i.POOL_ADDRESS_0_A))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())
		Expect(poolAccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
		Expect(poolAccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join pool again with different pool address after staker has left pool", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_A,
			Amount:        100 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.GetLeavePoolTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		_, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)
		Expect(active).To(BeFalse())

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:       i.STAKER_0,
			PoolId:        0,
			PoolAddress:   i.POOL_ADDRESS_0_B,
			Amount:        50 * i.KYVE,
			Commission:    math.LegacyMustNewDecFromStr("0.2"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterPoolAddress0 := s.GetBalanceFromAddress(i.POOL_ADDRESS_0_A)
		balanceAfterPoolAddress1 := s.GetBalanceFromAddress(i.POOL_ADDRESS_0_B)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(150 * i.KYVE))
		Expect(balanceAfterPoolAddress0 - initialBalancePoolAddress0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterPoolAddress1 - initialBalancePoolAddress1).To(Equal(50 * i.KYVE))

		poolAccountsOfStaker := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(poolAccountsOfStaker).To(HaveLen(1))

		poolAccount, active := s.App().StakersKeeper.GetPoolAccount(s.Ctx(), i.STAKER_0, 0)

		Expect(active).To(BeTrue())

		Expect(poolAccount.Staker).To(Equal(i.STAKER_0))
		Expect(poolAccount.PoolId).To(BeZero())
		Expect(poolAccount.PoolAddress).To(Equal(i.POOL_ADDRESS_0_B))
		Expect(poolAccount.Points).To(BeZero())
		Expect(poolAccount.IsLeaving).To(BeFalse())
		Expect(poolAccount.Commission).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
		Expect(poolAccount.StakeFraction).To(Equal(math.LegacyMustNewDecFromStr("1")))

		poolAccountsOfPool := s.App().StakersKeeper.GetAllPoolAccountsOfPool(s.Ctx(), 0)

		Expect(poolAccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().StakersKeeper.GetTotalStakeOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().StakersKeeper.GetValidatorPoolStake(s.Ctx(), i.STAKER_0, 0)).To(Equal(totalStakeOfPool))

		Expect(s.App().StakersKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})
})
