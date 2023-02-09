package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateMetadata = "update_metadata"

var _ sdk.Msg = &MsgUpdateMetadata{}

func (msg *MsgUpdateMetadata) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMetadata) Type() string {
	return TypeMsgUpdateMetadata
}

func (msg *MsgUpdateMetadata) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMetadata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMetadata) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Logo) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Logo), 255)
	}

	if len(msg.Website) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Website), 255)
	}

	if len(msg.Moniker) > 255 {
		return sdkerrors.Wrapf(sdkerrors.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Moniker), 255)
	}

	return nil
}
