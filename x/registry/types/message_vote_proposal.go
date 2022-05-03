package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgVoteProposal = "vote_proposal"

var _ sdk.Msg = &MsgVoteProposal{}

func NewMsgVoteProposal(creator string, id uint64, bundleId string, vote uint64) *MsgVoteProposal {
	return &MsgVoteProposal{
		Creator:  creator,
		Id:       id,
		BundleId: bundleId,
		Vote:  vote,
	}
}

func (msg *MsgVoteProposal) Route() string {
	return RouterKey
}

func (msg *MsgVoteProposal) Type() string {
	return TypeMsgVoteProposal
}

func (msg *MsgVoteProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
