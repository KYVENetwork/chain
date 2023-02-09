package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgVoteBundleProposal = "vote_bundle_proposal"

var _ sdk.Msg = &MsgVoteBundleProposal{}

func NewMsgVoteBundleProposal(creator string, staker string, poolId uint64, storageId string, vote VoteType) *MsgVoteBundleProposal {
	return &MsgVoteBundleProposal{
		Creator:   creator,
		Staker:    staker,
		PoolId:    poolId,
		StorageId: storageId,
		Vote:      vote,
	}
}

func (msg *MsgVoteBundleProposal) Route() string {
	return RouterKey
}

func (msg *MsgVoteBundleProposal) Type() string {
	return TypeMsgVoteBundleProposal
}

func (msg *MsgVoteBundleProposal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVoteBundleProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVoteBundleProposal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
