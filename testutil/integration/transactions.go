package integration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/gomega"
)

func (suite *KeeperTestSuite) RunTx(msg sdk.Msg) (*sdk.Result, error) {
	ctx, commit := suite.ctx.CacheContext()
	handler := suite.App().MsgServiceRouter().Handler(msg)

	res, err := handler(ctx, msg)
	if err != nil {
		return nil, err
	}

	commit()
	return res, nil
}

func (suite *KeeperTestSuite) RunTxError(msg sdk.Msg) error {
	_, err := suite.RunTx(msg)
	Expect(err).To(HaveOccurred())

	return err
}

func (suite *KeeperTestSuite) RunTxSuccess(msg sdk.Msg) *sdk.Result {
	result, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())

	return result
}

func (suite *KeeperTestSuite) RunTxGovSuccess(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxGovError(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).To(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxPoolSuccess(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxStakersSuccess(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxStakersError(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).To(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxBundlesSuccess(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxBundlesError(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).To(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxTeamSuccess(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxTeamError(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).To(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxFundersSuccess(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).NotTo(HaveOccurred())
}

func (suite *KeeperTestSuite) RunTxFundersError(msg sdk.Msg) {
	_, err := suite.RunTx(msg)
	Expect(err).To(HaveOccurred())
}
