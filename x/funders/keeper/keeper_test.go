package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KYVENetwork/chain/app/upgrades/v1_5"
	fundersOld "github.com/KYVENetwork/chain/app/upgrades/v1_5/v1_4_types/funders"
	i "github.com/KYVENetwork/chain/testutil/integration"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/x/funders/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFundersKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.ModuleName))
}

var _ = Describe("funders migrations", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		s = i.NewCleanChain()
	})

	It("v1.5", func() {
		// this panics
		// storeService := runtime.NewKVStoreService(storeTypes.NewKVStoreKey(types.StoreKey))
		// store := runtime.KVStoreAdapter(storeService.OpenKVStore(s.Ctx()))
		// fmt.Println(store.Get([]byte("test")))

		// ARRANGE
		storeKey, _ := v1_5.GetStoreKey(s.App().GetStoreKeys(), types.ModuleName)
		fundersOld.SetParams(s.Ctx(), s.App().AppCodec(), storeKey, fundersOld.Params{
			MinFundingAmount:          1,
			MinFundingAmountPerBundle: 2,
			MinFundingMultiple:        3,
		})

		// ACT
		v1_5.MigrateFundersModule(s.Ctx(), s.App().AppCodec(), s.App().GetStoreKeys(), s.App().FundersKeeper)

		// ASSERT
		params := s.App().FundersKeeper.GetParams(s.Ctx())

		Expect(params.MinFundingMultiple).To(Equal(uint64(3)))
		Expect(params.CoinWhitelist).To(HaveLen(1))
		Expect(params.CoinWhitelist[0].CoinDenom).To(Equal(globalTypes.Denom))
		Expect(params.CoinWhitelist[0].MinFundingAmount).To(Equal(uint64(1)))
		Expect(params.CoinWhitelist[0].MinFundingAmountPerBundle).To(Equal(uint64(2)))
	})
})
