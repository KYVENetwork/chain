package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgCreateTeamVestingAccount{}
	_ sdk.Msg            = &MsgCreateTeamVestingAccount{}
)

func (msg *MsgCreateTeamVestingAccount) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTeamVestingAccount) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTeamVestingAccount) Route() string {
	return RouterKey
}

func (msg *MsgCreateTeamVestingAccount) Type() string {
	return "kyve/team/MsgCreateTeamVestingAccount"
}

func (msg *MsgCreateTeamVestingAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
