package keeper_test

import (
	"fmt"
	"github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	. "github.com/onsi/ginkgo"
	"testing"
)

func TestAbs(t *testing.T) {
	s = new(KeeperTestSuite)
	s.SetupTest()

	//// Create and fund periodic vesting account
	//vestingStart := s.ctx.BlockTime()
	//baseAccount := authtypes.NewBaseAccountWithAddress(addr)
	//funder := sdk.AccAddress(types.ModuleName)

	ctx := sdk.WrapSDKContext(s.ctx)

	resp, err := keeper.NewMsgServerImpl(s.app.RegistryKeeper).FundPool(ctx, &types.MsgFundPool{
		Creator: "ich",
		Id:      0,
		Amount:  0,
	})

	fmt.Printf("%v\n%b\n", resp, err)

	//s.Require().Equal(vestingAmtTotal, unvested)
	s.Require().True(err == nil)
}

var _ = Describe("Books", func() {
	// Monthly vesting period
	stakeDenom := stakingtypes.DefaultParams().BondDenom
	amt := sdk.NewInt(1)
	vestingLength := int64(60 * 60 * 24 * 30) // in seconds
	vestingAmt := sdk.NewCoins(sdk.NewCoin(stakeDenom, amt))

	_ = vestingLength
	_ = vestingAmt

	BeforeEach(func() {
		s.SetupTest()

		//// Create and fund periodic vesting account
		//vestingStart := s.ctx.BlockTime()
		//baseAccount := authtypes.NewBaseAccountWithAddress(addr)
		//funder := sdk.AccAddress(types.ModuleName)

		resp, err := keeper.NewMsgServerImpl(s.app.RegistryKeeper).FundPool(s.ctx.Context(), &types.MsgFundPool{
			Creator: "ich",
			Id:      0,
			Amount:  0,
		})

		fmt.Printf("%v\n%b\n", resp, err)

		//s.Require().Equal(vestingAmtTotal, unvested)
		s.Require().True(err == nil)
	})

	Context("my first test", func() {
		It("It test", func() {
			fmt.Printf("Hello World\n")
		})
	})

	//Context("after first vesting period and before lockup", func() {
	//	BeforeEach(func() {
	//		// Surpass cliff but not lockup duration
	//		cliffDuration := time.Duration(cliffLength)
	//		s.CommitAfter(cliffDuration * time.Second)
	//
	//		// Check if some, but not all tokens are vested
	//		vested = clawbackAccount.GetVestedOnly(s.ctx.BlockTime())
	//		expVested := sdk.NewCoins(sdk.NewCoin(stakeDenom, amt.Mul(sdk.NewInt(cliff))))
	//		s.Require().NotEqual(vestingAmtTotal, vested)
	//		s.Require().Equal(expVested, vested)
	//	})
	//
	//	It("can delegate vested tokens", func() {
	//		err := delegate(clawbackAccount, vested.AmountOf(stakeDenom).Int64())
	//		Expect(err).To(BeNil())
	//	})
	//
	//	It("cannot delegate unvested tokens", func() {
	//		err := delegate(clawbackAccount, vestingAmtTotal.AmountOf(stakeDenom).Int64())
	//		Expect(err).ToNot(BeNil())
	//	})
	//
	//	It("cannot transfer vested tokens", func() {
	//		err := s.app.BankKeeper.SendCoins(
	//			s.ctx,
	//			addr,
	//			sdk.AccAddress(tests.GenerateAddress().Bytes()),
	//			vested,
	//		)
	//		Expect(err).ToNot(BeNil())
	//	})
	//
	//	It("cannot perform Ethereum tx", func() {
	//		err := performEthTx(clawbackAccount)
	//		Expect(err).ToNot(BeNil())
	//	})
	//})
	//
	//Context("after first vesting period and lockup", func() {
	//	BeforeEach(func() {
	//		// Surpass lockup duration
	//		lockupDuration := time.Duration(lockupLength)
	//		s.CommitAfter(lockupDuration * time.Second)
	//
	//		// Check if some, but not all tokens are vested
	//		unvested = clawbackAccount.GetUnvestedOnly(s.ctx.BlockTime())
	//		vested = clawbackAccount.GetVestedOnly(s.ctx.BlockTime())
	//		expVested := sdk.NewCoins(sdk.NewCoin(stakeDenom, amt.Mul(sdk.NewInt(lockup))))
	//		s.Require().NotEqual(vestingAmtTotal, vested)
	//		s.Require().Equal(expVested, vested)
	//	})
	//
	//	It("can delegate vested tokens", func() {
	//		err := delegate(clawbackAccount, vested.AmountOf(stakeDenom).Int64())
	//		Expect(err).To(BeNil())
	//	})
	//
	//	It("cannot delegate unvested tokens", func() {
	//		err := delegate(clawbackAccount, vestingAmtTotal.AmountOf(stakeDenom).Int64())
	//		Expect(err).ToNot(BeNil())
	//	})
	//
	//	It("can transfer vested tokens", func() {
	//		err := s.app.BankKeeper.SendCoins(
	//			s.ctx,
	//			addr,
	//			sdk.AccAddress(tests.GenerateAddress().Bytes()),
	//			vested,
	//		)
	//		Expect(err).To(BeNil())
	//	})
	//
	//	It("cannot transfer unvested tokens", func() {
	//		err := s.app.BankKeeper.SendCoins(
	//			s.ctx,
	//			addr,
	//			sdk.AccAddress(tests.GenerateAddress().Bytes()),
	//			unvested,
	//		)
	//		Expect(err).ToNot(BeNil())
	//	})
	//
	//	It("can perform ethereum tx", func() {
	//		err := performEthTx(clawbackAccount)
	//		Expect(err).To(BeNil())
	//	})
	//})
})
