package types

import (
	sdkErrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDelegate = "delegate"

var _ sdk.Msg = &MsgDelegate{}

func (msg *MsgDelegate) Route() string {
	return RouterKey
}

func (msg *MsgDelegate) Type() string {
	return TypeMsgDelegate
}

func (msg *MsgDelegate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDelegate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkErrors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Staker)
	if err != nil {
		return sdkErrors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid staker address (%s)", err)
	}
	return nil
}
