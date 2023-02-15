package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClawback = "clawback"

var _ sdk.Msg = &MsgClawback{}

func (msg *MsgClawback) Route() string {
	return RouterKey
}

func (msg *MsgClawback) Type() string {
	return TypeMsgClawback
}

func (msg *MsgClawback) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClawback) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClawback) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
