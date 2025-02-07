package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
)

/*

TEST CASES - msg_server_leave_pool.go

TODO
* Claim after period
* Check redistribution
* Check redistribution by weights

*/

var _ = Describe("logic_compliance_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var gov string
	var validator1 i.TestValidatorAddress

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()
		gov = s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

		// create pool
		msg := &pooltypes.MsgCreatePool{
			Authority:            gov,
			UploadInterval:       60,
			MaxBundleSize:        100,
			InflationShareWeight: math.LegacyZeroDec(),
			Binaries:             "{}",
		}
		s.RunTxPoolSuccess(msg)

		// create staker
		validator1 = s.CreateNewValidator("MyValidator-1", 1000*i.KYVE)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Claim after period is over", func() {
		// Arrange
		_ = validator1
	})
})
