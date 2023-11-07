package integration

import (
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) GetBalanceFromAddress(address string) uint64 {
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return 0
	}

	balance := suite.App().BankKeeper.GetBalance(suite.Ctx(), accAddress, globalTypes.Denom)

	return uint64(balance.Amount.Int64())
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
