package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateCommission{}

// GetSigners returns the expected signers for a MsgUpdateCommission message.
func (msg *MsgUpdateCommission) GetSigners() []sdk.AccAddress {
	validator, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{validator}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateCommission) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	if util.ValidatePercentage(msg.Commission) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid commission")
	}

	return nil
}
