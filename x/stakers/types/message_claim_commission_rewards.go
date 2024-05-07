package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgClaimCommissionRewards{}
	_ sdk.Msg            = &MsgClaimCommissionRewards{}
)

func (msg *MsgClaimCommissionRewards) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimCommissionRewards) GetSigners() []sdk.AccAddress {
	validator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{validator}
}

func (msg *MsgClaimCommissionRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimCommissionRewards) Type() string {
	return "kyve/stakers/MsgClaimCommissionRewards"
}

func (msg *MsgClaimCommissionRewards) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	if err := msg.Amount.Validate(); !msg.Amount.Empty() && err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount: %s", err)
	}

	return nil
}
