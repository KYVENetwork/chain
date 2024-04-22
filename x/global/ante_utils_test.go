package global_test

import (
	"cosmossdk.io/math"
	"cosmossdk.io/x/tx/signing"
	"github.com/KYVENetwork/chain/app"
	amino "github.com/cosmos/cosmos-sdk/codec"
	addressCodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gogoproto "github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/ibc-go/v8/testing/simapp/params"
	"google.golang.org/protobuf/proto"
)

// BuildEncodingConfig ...
func BuildEncodingConfig() params.EncodingConfig {
	cdc := amino.NewLegacyAmino()
	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: gogoproto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec:          addressCodec.NewBech32Codec(app.AccountAddressPrefix),
			ValidatorAddressCodec: addressCodec.NewBech32Codec(app.AccountAddressPrefix + "valoper"),
		},
	})
	if err != nil {
		panic(err)
	}
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
	txBuilder.SetFeePayer(sdk.MustAccAddressFromBech32(feePayer))

	msg := &banktypes.MsgSend{}
	_ = txBuilder.SetMsgs(msg)

	return txBuilder.GetTx()
}

// Invalid Transaction.
var _ sdk.Tx = &InvalidTx{}

type InvalidTx struct{}

func (t InvalidTx) GetMsgsV2() ([]proto.Message, error) {
	return []proto.Message{}, nil
}

func (InvalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (InvalidTx) ValidateBasic() error { return nil }

// AnteNextFn ...
func AnteNextFn(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}

// PostNextFn ...
func PostNextFn(ctx sdk.Context, _ sdk.Tx, _ bool, _ bool) (sdk.Context, error) {
	return ctx, nil
}
