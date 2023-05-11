package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgClaimCommissionRewards{}

// GetSigners returns the expected signers for a MsgUpdateCommission message.
func (msg *MsgClaimCommissionRewards) GetSigners() []sdk.AccAddress {
	validator, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{validator}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgClaimCommissionRewards) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	return nil
}
