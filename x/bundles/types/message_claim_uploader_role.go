package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgClaimUploaderRole{}
	_ sdk.Msg            = &MsgClaimUploaderRole{}
)

func NewMsgClaimUploaderRole(creator string, staker string, poolId uint64) *MsgClaimUploaderRole {
	return &MsgClaimUploaderRole{
		Creator: creator,
		Staker:  staker,
		PoolId:  poolId,
	}
}

func (msg *MsgClaimUploaderRole) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimUploaderRole) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimUploaderRole) Route() string {
	return RouterKey
}

func (msg *MsgClaimUploaderRole) Type() string {
	return "kyve/bundles/MsgClaimUploaderRole"
}

func (msg *MsgClaimUploaderRole) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
