package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgClaimAccountRewards{}
	_ sdk.Msg            = &MsgClaimAccountRewards{}
)

func (msg *MsgClaimAccountRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimAccountRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimAccountRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimAccountRewards) Type() string {
	return "kyve/team/MsgClaimAccountRewards"
}

func (msg *MsgClaimAccountRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
