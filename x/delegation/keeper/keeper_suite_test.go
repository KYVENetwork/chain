package keeper_test

import (
	"cosmossdk.io/math"
	"fmt"
	funderstypes "github.com/KYVENetwork/chain/x/funders/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDelegationKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.ModuleName))
}

func PayoutRewards(s *i.KeeperTestSuite, staker string, coins sdk.Coins) {
	fundingState, found := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)
	Expect(found).To(BeTrue())

	// TODO support all denominations once the funding module supports it
	// divide amount by number of active fundings so that total payout is equal to amount
	activeFundings := s.App().FundersKeeper.GetActiveFundings(s.Ctx(), fundingState)
	for _, funding := range activeFundings {
		funding.AmountsPerBundle = coins.QuoInt(math.NewInt(int64(len(activeFundings))))
		s.App().FundersKeeper.SetFunding(s.Ctx(), &funding)
	}

	payout, err := s.App().FundersKeeper.ChargeFundersOfPool(s.Ctx(), 0, pooltypes.ModuleName)
	Expect(err).To(BeNil())
	err = s.App().DelegationKeeper.PayoutRewards(s.Ctx(), staker, coins, pooltypes.ModuleName)
	Expect(err).NotTo(HaveOccurred())
	Expect(coins.String()).To(Equal(payout.String()))
}

func CreateFundedPool(s *i.KeeperTestSuite) {
	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	msg := &pooltypes.MsgCreatePool{
		Authority:            gov,
		Name:                 "PoolTest",
		Runtime:              "@kyve/test",
		Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
		Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
		StartKey:             "0",
		UploadInterval:       60,
		InflationShareWeight: 10_000,
		MinDelegation:        100 * i.KYVE,
		MaxBundleSize:        100,
		Version:              "0.0.0",
		Binaries:             "{}",
		StorageProviderId:    2,
		CompressionId:        1,
	}
	s.RunTxPoolSuccess(msg)

	s.CommitAfterSeconds(7)

	s.RunTxFundersSuccess(&funderstypes.MsgCreateFunder{
		Creator: i.ALICE,
		Moniker: "Alice",
	})

	s.RunTxPoolSuccess(&funderstypes.MsgFundPool{
		Creator:          i.ALICE,
		PoolId:           0,
		Amounts:          sdk.NewCoins(sdk.NewInt64Coin(i.KYVE_DENOM, 100*i.T_KYVE)),
		AmountsPerBundle: sdk.NewCoins(sdk.NewInt64Coin(i.KYVE_DENOM, 1*i.T_KYVE)),
	})

	s.CommitAfterSeconds(7)

	fundingState, _ := s.App().FundersKeeper.GetFundingState(s.Ctx(), 0)

	Expect(s.App().FundersKeeper.GetTotalActiveFunding(s.Ctx(), fundingState.PoolId)[0].Amount.Uint64()).To(Equal(100 * i.KYVE))
}

func CheckAndContinueChainForOneMonth(s *i.KeeperTestSuite) {
	s.PerformValidityChecks()

	for d := 0; d < 31; d++ {
		s.CommitAfterSeconds(60 * 60 * 24)
		s.PerformValidityChecks()
	}
}
