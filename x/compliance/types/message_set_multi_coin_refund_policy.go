package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgSetMultiCoinRewardsRefundPolicy{}
	_ sdk.Msg            = &MsgSetMultiCoinRewardsRefundPolicy{}
)

func (msg *MsgSetMultiCoinRewardsRefundPolicy) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetMultiCoinRewardsRefundPolicy) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgSetMultiCoinRewardsRefundPolicy) Route() string {
	return RouterKey
}

func (msg *MsgSetMultiCoinRewardsRefundPolicy) Type() string {
	return "kyve/stakers/MsgSetMultiCoinRewardsRefundPolicy"
}

func (msg *MsgSetMultiCoinRewardsRefundPolicy) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if _, err := ParseMultiCoinComplianceMap(*msg.Policy); err != nil {
		return err
	}

	return nil
}
