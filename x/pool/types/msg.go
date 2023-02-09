package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgUpdatePool{}
	_ sdk.Msg = &MsgDisablePool{}
	_ sdk.Msg = &MsgEnablePool{}
	_ sdk.Msg = &MsgScheduleRuntimeUpgrade{}
	_ sdk.Msg = &MsgCancelRuntimeUpgrade{}
)

// GetSigners returns the expected signers for a MsgCreatePool message.
func (msg *MsgCreatePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgCreatePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return nil
}

// GetSigners returns the expected signers for a MsgUpdatePool message.
func (msg *MsgUpdatePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdatePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return nil
}

// GetSigners returns the expected signers for a MsgDisablePool message.
func (msg *MsgDisablePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgDisablePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return nil
}

// GetSigners returns the expected signers for a MsgEnablePool message.
func (msg *MsgEnablePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgEnablePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return nil
}

// GetSigners returns the expected signers for a MsgScheduleRuntimeUpgrade message.
func (msg *MsgScheduleRuntimeUpgrade) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgScheduleRuntimeUpgrade) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return nil
}

// GetSigners returns the expected signers for a MsgCancelRuntimeUpgrade message.
func (msg *MsgCancelRuntimeUpgrade) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgCancelRuntimeUpgrade) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	return nil
}
