package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimAuthorityRewards = "claim_authority_rewards"

var _ sdk.Msg = &MsgClaimAuthorityRewards{}

func (msg *MsgClaimAuthorityRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimAuthorityRewards) Type() string {
	return TypeMsgClaimAuthorityRewards
}

func (msg *MsgClaimAuthorityRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimAuthorityRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimAuthorityRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
