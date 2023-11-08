package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	querytypes "github.com/KYVENetwork/chain/x/query/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - grpc_account_redelegation.go

* Call can validate if pool does not exist

*/

var _ = Describe("grpc_account_redelegation.go", Ordered, func() {
	s := i.NewCleanChain()

	redelegationCooldown := s.App().DelegationKeeper.GetRedelegationCooldown(s.Ctx())
	redelegationMaxAmount := s.App().DelegationKeeper.GetRedelegationMaxAmount(s.Ctx())

	BeforeEach(func() {
		s = i.NewCleanChain()

		// create 2 pools
		gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
		msg := &pooltypes.MsgCreatePool{
			Authority:      gov,
			MinDelegation:  200 * i.KYVE,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Binaries:       "{}",
		}
		s.RunTxPoolSuccess(msg)
		s.RunTxPoolSuccess(msg)

		// disable second pool
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		pool.Disabled = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
			Amount:     0,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.ALICE,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.STAKER_1,
			Amount:  50 * i.KYVE,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Single redelegation", func() {
		// ACT
		s.RunTxDelegatorSuccess(&delegationtypes.MsgRedelegate{
			Creator:    i.ALICE,
			FromStaker: i.STAKER_0,
			ToStaker:   i.STAKER_1,
			Amount:     10 * i.KYVE,
		})

		// ASSERT
		res, err := s.App().QueryKeeper.AccountRedelegation(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryAccountRedelegationRequest{Address: i.ALICE})
		Expect(err).To(BeNil())

		Expect(res.AvailableSlots).To(Equal(uint64(4)))
		Expect(res.RedelegationCooldownEntries).To(HaveLen(1))
		Expect(res.RedelegationCooldownEntries[0].CreationDate).To(Equal(uint64(s.Ctx().BlockTime().Unix())))
		Expect(res.RedelegationCooldownEntries[0].FinishDate).To(Equal(redelegationCooldown + uint64(s.Ctx().BlockTime().Unix())))
	})

	It("Await single redelegation", func() {
		// ACT
		s.RunTxDelegatorSuccess(&delegationtypes.MsgRedelegate{
			Creator:    i.ALICE,
			FromStaker: i.STAKER_0,
			ToStaker:   i.STAKER_1,
			Amount:     10 * i.KYVE,
		})
		s.CommitAfterSeconds(redelegationCooldown + 1)

		// Assert

		res, err := s.App().QueryKeeper.AccountRedelegation(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryAccountRedelegationRequest{Address: i.ALICE})
		Expect(err).To(BeNil())

		Expect(res.AvailableSlots).To(Equal(uint64(5)))
		Expect(res.RedelegationCooldownEntries).To(HaveLen(0))
	})

	It("Exhaust all redelegation", func() {
		// Arrange
		redelegationMsg := &delegationtypes.MsgRedelegate{
			Creator:    i.ALICE,
			FromStaker: i.STAKER_0,
			ToStaker:   i.STAKER_1,
			Amount:     10 * i.KYVE,
		}

		// ACT
		for i := uint64(0); i < redelegationMaxAmount; i++ {
			s.RunTxDelegatorSuccess(redelegationMsg)
			s.CommitAfterSeconds(1)
		}
		// Assert

		res, err := s.App().QueryKeeper.AccountRedelegation(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryAccountRedelegationRequest{Address: i.ALICE})
		Expect(err).To(BeNil())

		Expect(res.AvailableSlots).To(Equal(uint64(0)))
		Expect(res.RedelegationCooldownEntries).To(HaveLen(5))

		for i := uint64(0); i < redelegationMaxAmount; i++ {
			Expect(res.RedelegationCooldownEntries[i].CreationDate).To(Equal(uint64(s.Ctx().BlockTime().Unix()) - redelegationMaxAmount + i))
		}
	})
})
