package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgToggleMultiCoinRewards{}
	_ sdk.Msg            = &MsgToggleMultiCoinRewards{}
)

func (msg *MsgToggleMultiCoinRewards) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgToggleMultiCoinRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgToggleMultiCoinRewards) Route() string {
	return RouterKey
}

func (msg *MsgToggleMultiCoinRewards) Type() string {
	return "kyve/multi_coin_rewards/MsgToggleMultiCoinRewards"
}

func (msg *MsgToggleMultiCoinRewards) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	return nil
}
