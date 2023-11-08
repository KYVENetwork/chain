package keeper_test

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
)

/*

TEST CASES - msg_server_update_funder.go

* Try to update a funder that does not exist
* Try to update a funder with empty moniker
* Update a funder with empty values except moniker
* Update a funder with all values set
*/

var _ = Describe("msg_server_update_funder.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		moniker, identity, website, contact, description := "AliceMoniker", "identity", "website", "contact", "description"
		s.RunTxFundersSuccess(&types.MsgCreateFunder{
			Creator:     i.ALICE,
			Moniker:     moniker,
			Identity:    identity,
			Website:     website,
			Contact:     contact,
			Description: description,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Try to update a funder that does not exist", func() {
		// ASSERT
		s.RunTxFundersError(&types.MsgUpdateFunder{
			Creator: i.BOB,
			Moniker: "moniker",
		})
	})

	It("Try to update a funder with empty moniker", func() {
		// ASSERT
		s.RunTxFundersError(&types.MsgUpdateFunder{
			Creator: i.BOB,
			Moniker: "",
		})
	})

	It("Update a funder with empty values except moniker", func() {
		// ACT
		s.RunTxFundersSuccess(&types.MsgUpdateFunder{
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

	It("Update a funder with all values set", func() {
		// ACT
		moniker, identity, website, contact, description := "moniker", "identity", "website", "contact", "description"
		s.RunTxFundersSuccess(&types.MsgUpdateFunder{
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
})
