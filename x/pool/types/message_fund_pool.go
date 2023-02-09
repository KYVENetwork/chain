package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFundPool = "fund_pool"

var _ sdk.Msg = &MsgFundPool{}

func NewMsgFundPool(creator string, id uint64, amount uint64) *MsgFundPool {
	return &MsgFundPool{
		Creator: creator,
		Id:      id,
		Amount:  amount,
	}
}

func (msg *MsgFundPool) Route() string {
	return RouterKey
}

func (msg *MsgFundPool) Type() string {
	return TypeMsgFundPool
}

func (msg *MsgFundPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFundPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFundPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
