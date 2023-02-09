package types

import (
	sdkErrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRedelegate = "redelegate"

var _ sdk.Msg = &MsgRedelegate{}

func (msg *MsgRedelegate) Route() string {
	return RouterKey
}

func (msg *MsgRedelegate) Type() string {
	return TypeMsgRedelegate
}

func (msg *MsgRedelegate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRedelegate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkErrors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
