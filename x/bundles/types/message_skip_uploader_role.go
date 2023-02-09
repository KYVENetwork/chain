package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSkipUploaderRole = "skip_uploader_role"

var _ sdk.Msg = &MsgSkipUploaderRole{}

func NewMsgSkipUploaderRole(creator string, staker string, poolId uint64, fromIndex uint64) *MsgSkipUploaderRole {
	return &MsgSkipUploaderRole{
		Creator:   creator,
		Staker:    staker,
		PoolId:    poolId,
		FromIndex: fromIndex,
	}
}

func (msg *MsgSkipUploaderRole) Route() string {
	return RouterKey
}

func (msg *MsgSkipUploaderRole) Type() string {
	return TypeMsgSkipUploaderRole
}

func (msg *MsgSkipUploaderRole) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSkipUploaderRole) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSkipUploaderRole) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
