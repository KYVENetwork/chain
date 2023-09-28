package types

import (
	"encoding/json"

	"cosmossdk.io/errors"

	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgUpdatePool{}
	_ sdk.Msg = &MsgDisablePool{}
	_ sdk.Msg = &MsgEnablePool{}
	_ sdk.Msg = &MsgScheduleRuntimeUpgrade{}
	_ sdk.Msg = &MsgCancelRuntimeUpgrade{}
	_ sdk.Msg = &MsgUpdateParams{}
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

	if err := util.ValidatePositiveNumber(msg.UploadInterval); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid upload interval")
	}

	if err := util.ValidateNumber(msg.InflationShareWeight); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid inflation share weight")
	}

	if err := util.ValidateNumber(msg.MinDelegation); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid minimum delegation")
	}

	if err := util.ValidatePositiveNumber(msg.MaxBundleSize); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid max bundle size")
	}

	return nil
}

// GetSigners returns the expected signers for a MsgUpdatePool message.
func (msg *MsgUpdatePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// PoolUpdate ...
type PoolUpdate struct {
	Name                 *string
	Runtime              *string
	Logo                 *string
	Config               *string
	UploadInterval       *uint64
	InflationShareWeight *uint64
	MinDelegation        *uint64
	MaxBundleSize        *uint64
	StorageProviderId    *uint32
	CompressionId        *uint32
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdatePool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	var payload PoolUpdate
	if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
		return err
	}

	if payload.UploadInterval != nil {
		if err := util.ValidatePositiveNumber(*payload.UploadInterval); err != nil {
			return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid upload interval")
		}
	}

	if payload.InflationShareWeight != nil {
		if err := util.ValidateNumber(*payload.InflationShareWeight); err != nil {
			return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid inflation share weight")
		}
	}

	if payload.MinDelegation != nil {
		if err := util.ValidateNumber(*payload.MinDelegation); err != nil {
			return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid minimum delegation")
		}
	}

	if payload.MaxBundleSize != nil {
		if err := util.ValidatePositiveNumber(*payload.MaxBundleSize); err != nil {
			return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid max bundle size")
		}
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

// GetSigners returns the expected signers for a MsgCancelRuntimeUpgrade message.
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrap(err, "invalid authority address")
	}

	params := DefaultParams()
	if err := json.Unmarshal([]byte(msg.Payload), &params); err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	return nil
}
