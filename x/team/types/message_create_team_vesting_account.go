package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateTeamVestingAccount = "create_team_vesting_account"

var _ sdk.Msg = &MsgCreateTeamVestingAccount{}

func (msg *MsgCreateTeamVestingAccount) Route() string {
	return RouterKey
}

func (msg *MsgCreateTeamVestingAccount) Type() string {
	return TypeMsgCreateTeamVestingAccount
}

func (msg *MsgCreateTeamVestingAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTeamVestingAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTeamVestingAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
