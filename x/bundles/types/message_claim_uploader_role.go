package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimUploaderRole = "claim_uploader_role"

var _ sdk.Msg = &MsgClaimUploaderRole{}

func NewMsgClaimUploaderRole(creator string, staker string, poolId uint64) *MsgClaimUploaderRole {
	return &MsgClaimUploaderRole{
		Creator: creator,
		Staker:  staker,
		PoolId:  poolId,
	}
}

func (msg *MsgClaimUploaderRole) Route() string {
	return RouterKey
}

func (msg *MsgClaimUploaderRole) Type() string {
	return TypeMsgClaimUploaderRole
}

func (msg *MsgClaimUploaderRole) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimUploaderRole) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimUploaderRole) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
