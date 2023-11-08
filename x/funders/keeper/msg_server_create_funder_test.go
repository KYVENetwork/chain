package keeper_test

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
)

/*

TEST CASES - msg_server_create_funder.go

* Try to create a funder with empty values
* Create a funder with empty values except moniker
* Create a funder with all values set
* Try to create a funder that already exists
* Create two funders
*/

var _ = Describe("msg_server_create_funder.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Try to create a funder with empty values", func() {
		// ASSERT
		s.RunTxFundersError(&types.MsgCreateFunder{
			Creator: i.ALICE,
		})
	})

	It("Create a funder with empty values except moniker", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker",
		})

		// ASSERT
		funder, found := s.App().FundersKeeper.GetFunder(s.Ctx(), i.ALICE)
		Expect(found).To(BeTrue())
		Expect(funder.Address).To(Equal(i.ALICE))
		Expect(funder.Moniker).To(Equal("moniker"))
		Expect(funder.Identity).To(BeEmpty())
		Expect(funder.Website).To(BeEmpty())
		Expect(funder.Contact).To(BeEmpty())
		Expect(funder.Description).To(BeEmpty())
	})

	It("Create a funder with all values set", func() {
		// ACT
		moniker, identity, website, contact, description := "moniker", "identity", "website", "contact", "description"
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator:     i.ALICE,
			Moniker:     moniker,
			Identity:    identity,
			Website:     website,
			Contact:     contact,
			Description: description,
		})

		// ASSERT
		funder, found := s.App().FundersKeeper.GetFunder(s.Ctx(), i.ALICE)
		Expect(found).To(BeTrue())
		Expect(funder.Address).To(Equal(i.ALICE))
		Expect(funder.Moniker).To(Equal(moniker))
		Expect(funder.Identity).To(Equal(identity))
		Expect(funder.Website).To(Equal(website))
		Expect(funder.Contact).To(Equal(contact))
		Expect(funder.Description).To(Equal(description))
	})

	It("Try to create a funder that already exists", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker 1",
		})

		// ACT
		_, err := s.RunTx(&types.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker 2",
		})
		Expect(err.Error()).To(Equal("funder with address kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd already exists: invalid request"))
	})

	It("Create two funders", func() {
		// ARRANGE
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
			Moniker: "moniker 1",
		})

		// ACT
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator: i.BOB,
			Moniker: "moniker 2",
		})
	})
})
