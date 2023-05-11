package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgCreateStaker{}
	_ sdk.Msg            = &MsgCreateStaker{}
)

func (msg *MsgCreateStaker) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateStaker) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateStaker) Route() string {
	return RouterKey
}

func (msg *MsgCreateStaker) Type() string {
	return "kyve/stakers/MsgCreateStaker"
}

func (msg *MsgCreateStaker) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if util.ValidateNumber(msg.Amount) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount")
	}

	if msg.Commission.IsNil() {
		msg.Commission = DefaultCommission
	}
	if util.ValidatePercentage(msg.Commission) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid commission")
	}

	return nil
}
