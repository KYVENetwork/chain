package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFundPool{}

func (msg *MsgFundPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFundPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgFundPool) Route() string {
	return RouterKey
}

func (msg *MsgFundPool) Type() string {
	return "kyve/funders/MsgFundPool"
}

func (msg *MsgFundPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if util.ValidateNumber(msg.PoolId) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid pool id")
	}

	if util.ValidateNumber(msg.Amount) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount")
	}

	if util.ValidateNumber(msg.AmountPerBundle) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount per bundle")
	}

	return nil
}
