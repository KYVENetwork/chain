package types

import (
	sdkErrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUndelegate = "undelegate"

var _ sdk.Msg = &MsgUndelegate{}

func (msg *MsgUndelegate) Route() string {
	return RouterKey
}

func (msg *MsgUndelegate) Type() string {
	return TypeMsgUndelegate
}

func (msg *MsgUndelegate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUndelegate) ValidateBasic() error {
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
