package keeper_test

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
)

/*

TEST CASES - msg_server_create_funder.go

* Create a funder with empty values
* Create a funder with all values set
* Create a funder that already exists
* TODO: Create two funders with the same moniker
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

	It("Create a funder with empty values", func() {
		// ACT
		s.RunTxPoolSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
		})

		// ASSERT
		funder, found := s.App().FundersKeeper.GetFunder(s.Ctx(), i.ALICE)
		Expect(found).To(BeTrue())
		Expect(funder.Address).To(Equal(i.ALICE))
		Expect(funder.Moniker).To(BeEmpty())
		Expect(funder.Identity).To(BeEmpty())
		Expect(funder.Logo).To(BeEmpty())
		Expect(funder.Website).To(BeEmpty())
		Expect(funder.Contact).To(BeEmpty())
		Expect(funder.Description).To(BeEmpty())
	})

	It("Create a funder with all values set", func() {
		// ACT
		moniker, identity, logo, website, contact, description := "moniker", "identity", "logo", "website", "contact", "description"
		s.RunTxPoolSuccess(&types.MsgCreateFunder{
			Creator:     i.ALICE,
			Moniker:     moniker,
			Identity:    identity,
			Logo:        logo,
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
		Expect(funder.Logo).To(Equal(logo))
		Expect(funder.Website).To(Equal(website))
		Expect(funder.Contact).To(Equal(contact))
		Expect(funder.Description).To(Equal(description))
	})

	It("Create a funder with empty values", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&types.MsgCreateFunder{
			Creator: i.ALICE,
		})

		// ACT
		s.RunTxPoolError(&types.MsgCreateFunder{
			Creator: i.ALICE,
		})
	})
})
