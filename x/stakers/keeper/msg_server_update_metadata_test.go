package keeper_test

import (
	"cosmossdk.io/errors"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
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
		Expect(staker.Identity).To(BeEmpty())
		Expect(staker.SecurityContact).To(BeEmpty())
		Expect(staker.Details).To(BeEmpty())
	})

	It("Update metadata with real values of a newly created staker", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateMetadata{
			Creator:         i.STAKER_0,
			Moniker:         "KYVE Node Runner",
			Website:         "https://kyve.network",
			Identity:        "7CD454E228C8F227",
			SecurityContact: "security@kyve.network",
			Details:         "KYVE Protocol Node",
		})

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(staker.Moniker).To(Equal("KYVE Node Runner"))
		Expect(staker.Website).To(Equal("https://kyve.network"))
		Expect(staker.Identity).To(Equal("7CD454E228C8F227"))
		Expect(staker.SecurityContact).To(Equal("security@kyve.network"))
		Expect(staker.Details).To(Equal("KYVE Protocol Node"))
	})

	It("Reset metadata to empty values", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateMetadata{
			Creator:         i.STAKER_0,
			Moniker:         "KYVE Node Runner",
			Website:         "https://kyve.network",
			Identity:        "7CD454E228C8F227",
			SecurityContact: "security@kyve.network",
			Details:         "KYVE Protocol Node",
		})

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateMetadata{
			Creator:         i.STAKER_0,
			Moniker:         "",
			Website:         "",
			Identity:        "",
			SecurityContact: "",
			Details:         "",
		})

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)

		Expect(staker.Moniker).To(BeEmpty())
		Expect(staker.Website).To(BeEmpty())
		Expect(staker.Identity).To(BeEmpty())
		Expect(staker.SecurityContact).To(BeEmpty())
		Expect(staker.Details).To(BeEmpty())
	})

	It("One below max length", func() {
		// ARRANGE
		var stringStillAllowed string
		for i := 0; i < 255; i++ {
			stringStillAllowed += "."
		}

		// ACT
		msg := stakerstypes.MsgUpdateMetadata{
			Creator:         i.STAKER_0,
			Moniker:         stringStillAllowed,
			Website:         stringStillAllowed,
			Identity:        "",
			SecurityContact: stringStillAllowed,
			Details:         stringStillAllowed,
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
			Creator:         i.STAKER_0,
			Moniker:         stringTooLong,
			Website:         stringTooLong,
			Identity:        "",
			SecurityContact: stringTooLong,
			Details:         stringTooLong,
		}
		err := msg.ValidateBasic()

		// ASSERT
		Expect(err).ToNot(BeNil())
	})

	It("Invalid Identity", func() {
		// ARRANGE
		var invalidIdentity = "7CD454E228C8F22H"

		// ACT
		msg := stakerstypes.MsgUpdateMetadata{
			Creator:  i.STAKER_0,
			Identity: invalidIdentity,
		}
		err := msg.ValidateBasic()

		// ASSERT
		Expect(err.Error()).To(Equal(errors.Wrapf(errorsTypes.ErrLogic, stakerstypes.ErrInvalidIdentityString.Error(), msg.Identity).Error()))
	})

	It("Identity with lower-case hex letters", func() {
		// ARRANGE
		var invalidIdentity = "7cd454e228c8f227"

		// ACT
		msg := stakerstypes.MsgUpdateMetadata{
			Creator:  i.STAKER_0,
			Identity: invalidIdentity,
		}
		err := msg.ValidateBasic()

		// ASSERT
		Expect(err).To(BeNil())
	})
})
