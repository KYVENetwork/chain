package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgDefundPool{}
	_ sdk.Msg            = &MsgDefundPool{}
)

func NewMsgDefundPool(creator string, id uint64, amount uint64) *MsgDefundPool {
	return &MsgDefundPool{
		Creator: creator,
		Id:      id,
		Amount:  amount,
	}
}

func (msg *MsgDefundPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDefundPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgDefundPool) Route() string {
	return RouterKey
}

func (msg *MsgDefundPool) Type() string {
	return "kyve/pool/MsgDefundPool"
}

func (msg *MsgDefundPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
