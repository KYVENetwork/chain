package keeper_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"

	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"

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
	err := s.App().BankKeeper.MintCoins(s.Ctx(), mintTypes.ModuleName, coins)
	Expect(err).NotTo(HaveOccurred())

	s.Commit()

	err = s.App().DelegationKeeper.PayoutRewards(s.Ctx(), staker, coins, mintTypes.ModuleName)
	Expect(err).NotTo(HaveOccurred())
}

func CreatePool(s *i.KeeperTestSuite) {
	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	msg := &pooltypes.MsgCreatePool{
		Authority:            gov,
		Name:                 "PoolTest",
		Runtime:              "@kyve/test",
		Logo:                 "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
		Config:               "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
		StartKey:             "0",
		UploadInterval:       60,
		InflationShareWeight: math.LegacyNewDec(10_000),
		MinDelegation:        100 * i.KYVE,
		MaxBundleSize:        100,
		Version:              "0.0.0",
		Binaries:             "{}",
		StorageProviderId:    2,
		CompressionId:        1,
	}
	s.RunTxPoolSuccess(msg)
}

func CheckAndContinueChainForOneMonth(s *i.KeeperTestSuite) {
	s.PerformValidityChecks()

	for d := 0; d < 31; d++ {
		s.CommitAfterSeconds(60 * 60 * 24)
		s.PerformValidityChecks()
	}
}
