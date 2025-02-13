package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

var (
	_ legacytx.LegacyMsg = &MsgSetMultiCoinRewardsDistributionPolicy{}
	_ sdk.Msg            = &MsgSetMultiCoinRewardsDistributionPolicy{}
)

func (msg *MsgSetMultiCoinRewardsDistributionPolicy) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetMultiCoinRewardsDistributionPolicy) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgSetMultiCoinRewardsDistributionPolicy) Route() string {
	return RouterKey
}

func (msg *MsgSetMultiCoinRewardsDistributionPolicy) Type() string {
	return "kyve/multi_coin_rewards/MsgSetMultiCoinRewardsDistributionPolicy"
}

func (msg *MsgSetMultiCoinRewardsDistributionPolicy) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	if msg.Policy == nil {
		return errors.Wrap(errorsTypes.ErrInvalidRequest, "policy cannot be nil")
	}

	if _, err := ParseAndNormalizeMultiCoinDistributionMap(*msg.Policy); err != nil {
		return err
	}

	return nil
}
