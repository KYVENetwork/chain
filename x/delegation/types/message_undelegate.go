package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgUndelegate{}
	_ sdk.Msg            = &MsgUndelegate{}
)

func (msg *MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUndelegate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgUndelegate) Route() string {
	return RouterKey
}

func (msg *MsgUndelegate) Type() string {
	return "kyve/delegation/MsgUndelegate"
}

func (msg *MsgUndelegate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Staker)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid staker address (%s)", err)
	}

	return nil
}
