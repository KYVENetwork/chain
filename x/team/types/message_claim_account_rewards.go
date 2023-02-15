package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimAccountRewards = "claim_account_rewards"

var _ sdk.Msg = &MsgClaimAccountRewards{}

func (msg *MsgClaimAccountRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimAccountRewards) Type() string {
	return TypeMsgClaimAccountRewards
}

func (msg *MsgClaimAccountRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimAccountRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimAccountRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
