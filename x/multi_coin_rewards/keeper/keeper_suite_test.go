package keeper_test

import (
	"fmt"
	"testing"

	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiCoinRewardsKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.ModuleName))
}

func payoutRewards(s *i.KeeperTestSuite, staker string, coins sdk.Coins) {
	err := s.App().BankKeeper.MintCoins(s.Ctx(), mintTypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}
	err = s.App().StakersKeeper.PayoutRewards(s.Ctx(), staker, coins, mintTypes.ModuleName)
	if err != nil {
		panic(err)
	}
}
