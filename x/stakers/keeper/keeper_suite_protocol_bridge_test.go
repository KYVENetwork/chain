package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_join_pool.go

* Consensus slash leads to removal from pool
* Consensus slash leads to removal from all pool
* Getting pushed out of the active set leads to removal from pool
* Unbonded validator can not join a pool

*/

var _ = Describe("keeper_suite_protocol_bridge.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	validator1 := s.CreateNewValidator("Staker-0", 1000*i.KYVE)

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

		validator1 = s.CreateNewValidator("Staker-0", 1000*i.KYVE)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Consensus slash leads to removal from pool", func() {
		// Arrange
		params, _ := s.App().SlashingKeeper.GetParams(s.Ctx())
		params.MinSignedPerWindow = math.LegacyMustNewDecFromStr("1")
		params.SignedBlocksWindow = 1
		_ = s.App().SlashingKeeper.SetParams(s.Ctx(), params)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       validator1.Address,
			PoolId:        0,
			PoolAddress:   validator1.PoolAccount[0],
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// Make validator not participate in block votes to have him kicked out
		s.AddAbciAbsentVote(validator1.ConsAccAddress)

		preBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(preBonded).To(HaveLen(2))

		// Act
		s.CommitAfterSeconds(1)
		s.CommitAfterSeconds(1)

		// Assert
		poolMembersCount := s.App().StakersKeeper.GetStakerCountOfPool(s.Ctx(), 0)
		Expect(poolMembersCount).To(Equal(uint64(0)))

		postBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(postBonded).To(HaveLen(1))

		poolAccounts := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), validator1.Address)
		Expect(poolAccounts).To(HaveLen(0))
	})

	It("Consensus slash leads to removal from all pool", func() {
		// Arrange
		params, _ := s.App().SlashingKeeper.GetParams(s.Ctx())
		params.MinSignedPerWindow = math.LegacyMustNewDecFromStr("1")
		params.SignedBlocksWindow = 1
		_ = s.App().SlashingKeeper.SetParams(s.Ctx(), params)

		for k := 1; k < 10; k++ {
			msg := &pooltypes.MsgCreatePool{
				Authority:            gov,
				UploadInterval:       60,
				MaxBundleSize:        100,
				InflationShareWeight: math.LegacyZeroDec(),
				Binaries:             "{}",
			}
			s.RunTxPoolSuccess(msg)

			s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
				Creator:       validator1.Address,
				PoolId:        uint64(k),
				PoolAddress:   validator1.PoolAccount[k],
				Commission:    math.LegacyMustNewDecFromStr("0.1"),
				StakeFraction: math.LegacyMustNewDecFromStr("1"),
			})
		}

		// Make validator not participate in block votes to have him kicked out
		s.AddAbciAbsentVote(validator1.ConsAccAddress)

		preBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(preBonded).To(HaveLen(2))

		poolAccounts := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), validator1.Address)
		Expect(poolAccounts).To(HaveLen(9))

		// Act
		s.CommitAfterSeconds(1)
		s.CommitAfterSeconds(1)

		// Assert
		for k := 1; k < 10; k++ {
			poolMembersCount := s.App().StakersKeeper.GetStakerCountOfPool(s.Ctx(), uint64(k))
			Expect(poolMembersCount).To(Equal(uint64(0)))
		}

		postBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(postBonded).To(HaveLen(1))

		poolAccounts = s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), validator1.Address)
		Expect(poolAccounts).To(HaveLen(0))
	})

	It("Getting pushed out of the active set leads to removal from pool", func() {
		// Arrange

		// There are 51 validator slots (1 is occupied by the default validator (1 $KYVE) and another one by validator1)
		for k := 0; k < 50; k++ {
			s.CreateNewValidator(fmt.Sprintf("val-%d", k), 10000*i.KYVE)
		}
		s.CommitAfterSeconds(1)
		s.CommitAfterSeconds(1)

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:       validator1.Address,
			PoolId:        0,
			PoolAddress:   validator1.PoolAccount[0],
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		preBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(preBonded).To(HaveLen(51))

		// Act

		// Create one more validator to kick out validator1
		s.CreateNewValidator(fmt.Sprintf("val-%d", 51), 10000*i.KYVE)
		s.CommitAfterSeconds(1)
		s.CommitAfterSeconds(1)

		// Assert
		poolMembersCount := s.App().StakersKeeper.GetStakerCountOfPool(s.Ctx(), 0)
		Expect(poolMembersCount).To(Equal(uint64(0)))

		postBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(postBonded).To(HaveLen(51))

		poolAccounts := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), validator1.Address)
		Expect(poolAccounts).To(HaveLen(0))
	})

	It("Unbonded validator can not join a pool", func() {
		// Arrange

		// There are 51 validator slots (1 is occupied by the default validator (1 $KYVE) and another one by validator1)
		// This will kick out validator1
		for k := 0; k < 51; k++ {
			s.CreateNewValidator(fmt.Sprintf("val-%d", k), 10000*i.KYVE)
		}
		s.CommitAfterSeconds(1)
		s.CommitAfterSeconds(1)

		preBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(preBonded).To(HaveLen(51))

		// Act
		_, err := s.RunTx(&stakertypes.MsgJoinPool{
			Creator:       validator1.Address,
			PoolId:        0,
			PoolAddress:   validator1.PoolAccount[0],
			Commission:    math.LegacyMustNewDecFromStr("0.1"),
			StakeFraction: math.LegacyMustNewDecFromStr("1"),
		})

		// Assert
		Expect(err.Error()).To(Equal("validator not in active set"))

		poolMembersCount := s.App().StakersKeeper.GetStakerCountOfPool(s.Ctx(), 0)
		Expect(poolMembersCount).To(Equal(uint64(0)))

		postBonded, _ := s.App().StakingKeeper.GetBondedValidatorsByPower(s.Ctx())
		Expect(postBonded).To(HaveLen(51))

		poolAccounts := s.App().StakersKeeper.GetPoolAccountsFromStaker(s.Ctx(), validator1.Address)
		Expect(poolAccounts).To(HaveLen(0))
	})
})
