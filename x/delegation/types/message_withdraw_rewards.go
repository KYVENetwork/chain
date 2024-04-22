package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgWithdrawRewards{}
	_ sdk.Msg            = &MsgWithdrawRewards{}
)

func (msg *MsgWithdrawRewards) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawRewards) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawRewards) Type() string {
	return "kyve/delegation/MsgWithdrawRewards"
}

func (msg *MsgWithdrawRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Staker)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid staker address (%s)", err)
	}

	return nil
}
