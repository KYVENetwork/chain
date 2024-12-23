package integration

import (
	"fmt"

	"cosmossdk.io/math"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) MustNewIntFromStr(amount string) math.Int {
	result, ok := math.NewIntFromString(amount)
	if !ok {
		panic(fmt.Sprintf("error parsing \"%s\" to math.Int", amount))
	}
	return result
}

func (suite *KeeperTestSuite) GetCoinsFromCommunityPool() sdk.Coins {
	pool, err := suite.App().DistributionKeeper.FeePool.Get(suite.Ctx())
	if err != nil {
		return sdk.NewCoins()
	}

	coins, _ := pool.CommunityPool.TruncateDecimal()
	return coins
}

func (suite *KeeperTestSuite) GetBalanceFromAddress(address string) uint64 {
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return 0
	}

	balance := suite.App().BankKeeper.GetBalance(suite.Ctx(), accAddress, globalTypes.Denom)

	return uint64(balance.Amount.Int64())
}

func (suite *KeeperTestSuite) GetCoinsFromAddress(address string) sdk.Coins {
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdk.NewCoins()
	}

	return suite.App().BankKeeper.GetAllBalances(suite.Ctx(), accAddress)
}

func (suite *KeeperTestSuite) GetBalanceFromPool(poolId uint64) uint64 {
	pool, found := suite.App().PoolKeeper.GetPool(suite.Ctx(), poolId)
	if !found {
		return 0
	}

	return uint64(suite.App().BankKeeper.GetBalance(suite.Ctx(), pool.GetPoolAccount(), globalTypes.Denom).Amount.Int64())
}

func (suite *KeeperTestSuite) GetBalanceFromModule(moduleName string) uint64 {
	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), moduleName).GetAddress()
	return suite.App().BankKeeper.GetBalance(suite.Ctx(), moduleAcc, globalTypes.Denom).Amount.Uint64()
}

func (suite *KeeperTestSuite) GetCoinsFromModule(moduleName string) sdk.Coins {
	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), moduleName).GetAddress()
	return suite.App().BankKeeper.GetAllBalances(suite.Ctx(), moduleAcc)
}

func (suite *KeeperTestSuite) GetNextUploader() (nextStaker string, nextValaddress string) {
	bundleProposal, _ := suite.App().BundlesKeeper.GetBundleProposal(suite.Ctx(), 0)

	switch bundleProposal.NextUploader {
	case STAKER_0:
		nextStaker = STAKER_0
		nextValaddress = VALADDRESS_0_A
	case STAKER_1:
		nextStaker = STAKER_1
		nextValaddress = VALADDRESS_1_A
	case STAKER_2:
		nextStaker = STAKER_2
		nextValaddress = VALADDRESS_2_A
	default:
		nextStaker = ""
		nextValaddress = ""
	}

	return
}

func (suite *KeeperTestSuite) SetMaxVotingPower(maxVotingPower string) {
	params := suite.App().PoolKeeper.GetParams(suite.Ctx())
	params.MaxVotingPowerPerPool = math.LegacyMustNewDecFromStr(maxVotingPower)
	suite.App().PoolKeeper.SetParams(suite.Ctx(), params)
}
