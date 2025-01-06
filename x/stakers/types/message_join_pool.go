package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgJoinPool{}
	_ sdk.Msg            = &MsgJoinPool{}
)

func (msg *MsgJoinPool) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgJoinPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgJoinPool) Route() string {
	return RouterKey
}

func (msg *MsgJoinPool) Type() string {
	return "kyve/stakers/MsgJoinPool"
}

func (msg *MsgJoinPool) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.PoolAddress); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid pool address: %s", err)
	}

	if util.ValidateNumber(msg.Amount) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid amount")
	}

	if util.ValidatePercentage(msg.Commission) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid commission")
	}

	if util.ValidatePercentage(msg.StakeFraction) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid stake fraction")
	}

	return nil
}
