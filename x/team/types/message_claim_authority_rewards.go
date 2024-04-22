package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgClaimAuthorityRewards{}
	_ sdk.Msg            = &MsgClaimAuthorityRewards{}
)

func (msg *MsgClaimAuthorityRewards) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimAuthorityRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimAuthorityRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimAuthorityRewards) Type() string {
	return "kyve/team/MsgClaimAuthorityRewards"
}

func (msg *MsgClaimAuthorityRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
