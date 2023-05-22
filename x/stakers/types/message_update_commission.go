package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgUpdateCommission{}
	_ sdk.Msg            = &MsgUpdateCommission{}
)

func (msg *MsgUpdateCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCommission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCommission) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCommission) Type() string {
	return "kyve/stakers/MsgUpdateCommission"
}

func (msg *MsgUpdateCommission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	if util.ValidatePercentage(msg.Commission) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid commission")
	}

	return nil
}
