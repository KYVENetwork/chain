package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - abci.go

* inactive pool should not receive inflation funds
* active pool should receive inflation funds
* pool should split inflation funds depending on inflation share weight
* pools with zero inflation share weight should receive nothing
* every pool has zero inflation share weight

*/

var _ = Describe("abci.go", Ordered, func() {
	s := i.NewCleanChain()
	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	BeforeEach(func() {
		s = i.NewCleanChain()

		s.App().PoolKeeper.SetParams(s.Ctx(), pooltypes.Params{
			ProtocolInflationShare:  sdk.MustNewDecFromStr("0.1"),
			PoolInflationPayoutRate: sdk.MustNewDecFromStr("0.1"),
		})

		// create clean pool for every test case
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 1_000_000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("inactive pool should not receive inflation funds", func() {
		// ARRANGE
		b1, b2 := uint64(0), uint64(0)

		for t := 0; t < 100; t++ {
			// ACT
			pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
			b1 = uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())
			s.Commit()
			b2 = uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

			// ASSERT
			Expect(b1).To(BeZero())
			Expect(b2).To(BeZero())
		}
	})

	It("active pool should receive inflation funds", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
		})

		b1, b2 := uint64(0), uint64(0)

		for t := 0; t < 100; t++ {
			// ACT
			pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
			b1 = uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())
			s.Commit()
			b2 = uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

			// ASSERT
			Expect(b1).To(BeNumerically("<", b2))
		}
	})

	It("pool should split inflation funds depending on inflation share weight", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 2_000_000,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0_B,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1_B,
		})

		// ACT
		s.Commit()

		pool1, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		b1 := uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool1.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

		pool2, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		b2 := uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool2.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

		// ASSERT
		Expect(b1 * 2).To(BeNumerically("~", b2-10, b2+10))
	})

	It("pools with zero inflation share weight should receive nothing", func() {
		// ARRANGE
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 0,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0_A,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1_A,
		})

		// ACT
		s.Commit()

		pool1, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		b1 := uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool1.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

		pool2, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		b2 := uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool2.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

		// ASSERT
		Expect(b1).To(BeNumerically(">", b2))
		Expect(b2).To(BeZero())
	})

	It("every pool has zero inflation share weight", func() {
		// ARRANGE
		pool1, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool1.InflationShareWeight = 0
		s.App().PoolKeeper.SetPool(s.Ctx(), pool1)

		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			Name:                 "PoolTest",
			Runtime:              "@kyve/test",
			Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:             "0",
			UploadInterval:       60,
			InflationShareWeight: 0,
			MinDelegation:        100 * i.KYVE,
			MaxBundleSize:        100,
			Version:              "0.0.0",
			Binaries:             "{}",
			StorageProviderId:    2,
			CompressionId:        1,
		}
		s.RunTxPoolSuccess(msg)

		// ACT
		s.Commit()

		pool1, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		b1 := uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool1.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

		pool2, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		b2 := uint64(s.App().BankKeeper.GetBalance(s.Ctx(), pool2.GetPoolAccount(), globalTypes.Denom).Amount.Int64())

		// ASSERT
		Expect(b1).To(BeZero())
		Expect(b2).To(BeZero())
	})
})
