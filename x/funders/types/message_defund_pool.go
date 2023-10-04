package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgDefundPool{}

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
	return "kyve/funders/MsgDefundPool"
}

func (msg *MsgDefundPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if util.ValidateNumber(msg.PoolId) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid pool id")
	}

	if util.ValidatePositiveNumber(msg.Amount) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount")
	}

	return nil
}
