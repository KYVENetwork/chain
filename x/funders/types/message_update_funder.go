package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateFunder{}

func (msg *MsgUpdateFunder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateFunder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateFunder) Route() string {
	return RouterKey
}

func (msg *MsgUpdateFunder) Type() string {
	return "kyve/funders/MsgUpdateFunder"
}

func (msg *MsgUpdateFunder) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}
	if msg.Moniker == "" {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "moniker cannot be empty")
	}

	return nil
}
