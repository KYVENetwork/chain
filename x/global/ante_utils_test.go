package global_test

import (
	"cosmossdk.io/math"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

// BuildEncodingConfig ...
func BuildEncodingConfig() params.EncodingConfig {
	cdc := amino.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	codec := amino.NewProtoCodec(interfaceRegistry)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             cdc,
	}

	return encodingConfig
}

// BuildTestTx ...
func BuildTestTx(gasPrice math.Int, denom string, feePayer string, encodingConfig params.EncodingConfig) sdk.FeeTx {
	txBuilder := encodingConfig.TxConfig.NewTxBuilder()

	gasLimit := uint64(200_000)
	txBuilder.SetGasLimit(gasLimit)

	fees := sdk.NewCoins(sdk.NewCoin(denom, gasPrice.MulRaw(int64(gasLimit))))
	txBuilder.SetFeeAmount(fees)

	msg := &TestMsg{Signers: []string{feePayer}}
	_ = txBuilder.SetMsgs(msg)

	return txBuilder.GetTx()
}

// Invalid Transaction.
var _ sdk.Tx = &InvalidTx{}

type InvalidTx struct{}

func (InvalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (InvalidTx) ValidateBasic() error { return nil }

// NextFn ...
func NextFn(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}

// Test Message.
var _ sdk.Msg = (*TestMsg)(nil)

type TestMsg struct {
	Signers []string
}

func (msg *TestMsg) Reset()               {}
func (msg *TestMsg) String() string       { return "" }
func (msg *TestMsg) ProtoMessage()        {}
func (msg *TestMsg) ValidateBasic() error { return nil }

func (msg *TestMsg) GetSigners() []sdk.AccAddress {
	var addrs []sdk.AccAddress

	for _, signer := range msg.Signers {
		addr := sdk.MustAccAddressFromBech32(signer)
		addrs = append(addrs, addr)
	}

	return addrs
}
