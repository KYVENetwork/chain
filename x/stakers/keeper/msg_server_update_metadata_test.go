package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_update_metadata.go

* Get the default metadata of a newly created staker
* Update metadata with real values of a newly created staker
* Reset metadata to empty values
* Exceed max length

*/

var _ = Describe("msg_server_update_metadata.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create staker
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Get the default metadata of a newly created staker", func() {
		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())
		Expect(staker.Logo).To(BeEmpty())
	})

	It("Update metadata with real values of a newly created staker", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateMetadata{
			Creator: i.STAKER_0,
			Moniker: "KYVE Node Runner",
			Website: "https://kyve.network",
			Logo:    "https://arweave.net/Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
		})

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(staker.Moniker).To(Equal("KYVE Node Runner"))
		Expect(staker.Website).To(Equal("https://kyve.network"))
		Expect(staker.Logo).To(Equal("https://arweave.net/Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU"))
	})

	It("Reset metadata to empty values", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateMetadata{
			Creator: i.STAKER_0,
			Moniker: "KYVE Node Runner",
			Website: "https://kyve.network",
			Logo:    "https://arweave.net/Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateMetadata{
			Creator: i.STAKER_0,
			Moniker: "",
			Website: "",
			Logo:    "",
		})

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())
		Expect(staker.Logo).To(BeEmpty())
	})

	It("One below max length", func() {
		// ARRANGE
		var stringStillAllowed string
		for i := 0; i < 255; i++ {
			stringStillAllowed += "."
		}

		// ACT
		msg := stakerstypes.MsgUpdateMetadata{
			Creator: i.STAKER_0,
			Moniker: stringStillAllowed,
			Website: stringStillAllowed,
			Logo:    stringStillAllowed,
		}
		err := msg.ValidateBasic()

		// ASSERT
		Expect(err).To(BeNil())
	})

	// stringTooLong := stringStillAllowed + "."
	It("Exceed max length", func() {
		// ARRANGE
		var stringTooLong string
		for i := 0; i < 256; i++ {
			stringTooLong += "."
		}

		// ACT
		msg := stakerstypes.MsgUpdateMetadata{
			Creator: i.STAKER_0,
			Moniker: stringTooLong,
			Website: stringTooLong,
			Logo:    stringTooLong,
		}
		err := msg.ValidateBasic()

		// ASSERT
		Expect(err).ToNot(BeNil())
	})
})
