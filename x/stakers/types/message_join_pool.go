package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgJoinPool{}

func (msg *MsgJoinPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgJoinPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Valaddress); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	if util.ValidateNumber(msg.Amount) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount")
	}

	return nil
}
