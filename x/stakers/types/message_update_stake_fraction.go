package types

import (
	"cosmossdk.io/errors"
	"github.com/KYVENetwork/chain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgUpdateStakeFraction{}
	_ sdk.Msg            = &MsgUpdateStakeFraction{}
)

func (msg *MsgUpdateStakeFraction) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateStakeFraction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateStakeFraction) Route() string {
	return RouterKey
}

func (msg *MsgUpdateStakeFraction) Type() string {
	return "kyve/stakers/MsgUpdateStakeFraction"
}

func (msg *MsgUpdateStakeFraction) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid validator address: %s", err)
	}

	if util.ValidatePercentage(msg.StakeFraction) != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidRequest, "invalid stake fraction")
	}

	return nil
}
