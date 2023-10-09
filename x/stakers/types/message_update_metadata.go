package types

import (
	"encoding/hex"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgUpdateMetadata{}
	_ sdk.Msg            = &MsgUpdateMetadata{}
)

func (msg *MsgUpdateMetadata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMetadata) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMetadata) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMetadata) Type() string {
	return "kyve/stakers/MsgUpdateMetadata"
}

func (msg *MsgUpdateMetadata) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Identity) > 0 {
		if hexBytes, identityErr := hex.DecodeString(msg.Identity); identityErr != nil || len(hexBytes) != 8 {
			return errors.Wrapf(errorsTypes.ErrLogic, ErrInvalidIdentityString.Error(), msg.Identity)
		}
	}

	if len(msg.Website) > 255 {
		return errors.Wrapf(errorsTypes.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Website), 255)
	}

	if len(msg.Moniker) > 255 {
		return errors.Wrapf(errorsTypes.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Moniker), 255)
	}

	if len(msg.SecurityContact) > 255 {
		return errors.Wrapf(errorsTypes.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Moniker), 255)
	}

	if len(msg.Details) > 255 {
		return errors.Wrapf(errorsTypes.ErrLogic, ErrStringMaxLengthExceeded.Error(), len(msg.Moniker), 255)
	}

	return nil
}
