package integration

import (
	"encoding/hex"
	mrand "math/rand"
	"time"

	"github.com/stretchr/testify/suite"

	abci "github.com/cometbft/cometbft/abci/types"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/app"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	"github.com/cometbft/cometbft/version"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	ALICE   = "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"
	BOB     = "kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq"
	CHARLIE = "kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s"
	DAVID   = "kyve1jxa7kp37jlm8hzwgc5qprquv9k7vawq79qhctt"

	STAKER_0       = "kyve1htgfatqevuvfzvl0sxp97ywteqhg5leha9emf4"
	VALADDRESS_0_A = "kyve1qnf86dkvvtpdukx30r3vajav7rdq8snktm90hm"
	VALADDRESS_0_B = "kyve10t8gnqjnem7tsu09erzswj3zm8599lsnex79rz"
	VALADDRESS_0_C = "kyve13ztkcm2pket6mrmxj8rrmwc6supw7aqakg3uu3"

	STAKER_1       = "kyve1gnr35rwn8rmflnlzs6nn5hhkmzzkxg9ap8xepw"
	VALADDRESS_1_A = "kyve1hpjgzljglmv00nstk3jvcw0zzq94nu0cuxv5ga"
	VALADDRESS_1_B = "kyve14runw9qkltpz2mcx3gsfmlqyyvdzkt3rq3w6fm"
	VALADDRESS_1_C = "kyve15w9m7zpq9ctsxsveqaqkp4uuvw98z5vct6s9g9"

	STAKER_2       = "kyve1xsemlxghgvusumhqzm2ztjw7dz9krvu3de54e2"
	VALADDRESS_2_A = "kyve1u0870dkae6ql63hxvy9y7g65c0y8csfh8allzl"
	VALADDRESS_2_B = "kyve16g3utghkvvlz53jk0fq96zwrhxmqfu36ue965q"
	VALADDRESS_2_C = "kyve18gjtzsn6jme3qsczj9q7wefymlkfu7ngyq5f9c"

	STAKER_3       = "kyve1ca7rzyrxfpdm7j8jgccq4rduuf4sxpq0dhmwm4"
	VALADDRESS_3_A = "kyve1d2clkfrw0r99ctgmkjvluzn6xm98yls06mnxv8"
	VALADDRESS_3_B = "kyve1f36cvde6jnygcrz2yas4acp0akn9cw7vp5ze0w"
	VALADDRESS_3_C = "kyve1gcnd8gya2ysfur6d6z4wpl9z54zadg7qzk8uyc"

	// To avoid giving burner permissions to a module for the tests
	BURNER = "kyve1ld23ktfwc9zstaq8aanwkkj8cf0ru6adtz59y5"
)

var (
	DUMMY    []string
	VALDUMMY []string
)

const (
	KYVE   = uint64(1_000_000_000)
	T_KYVE = int64(KYVE)
)

var (
	KYVE_DENOM = globalTypes.Denom
	A_DENOM    = "zcoin"
	B_DENOM    = "bcoin"
	C_DENOM    = "ccoin"
)

func KYVECoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(KYVE_DENOM, amount)
}

func KYVECoins(amount int64) sdk.Coins {
	return sdk.NewCoins(KYVECoin(amount))
}

func ACoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(A_DENOM, amount)
}

func ACoins(amount int64) sdk.Coins {
	return sdk.NewCoins(ACoin(amount))
}

func BCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(B_DENOM, amount)
}

func BCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(BCoin(amount))
}

func CCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(C_DENOM, amount)
}

func CCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(CCoin(amount))
}

func NewCleanChain() *KeeperTestSuite {
	s := KeeperTestSuite{}
	s.SetupTest(time.Now().Unix())
	s.initDummyAccounts()
	return &s
}

func NewCleanChainAtTime(startTime int64) *KeeperTestSuite {
	s := KeeperTestSuite{}
	s.SetupTest(startTime)
	s.initDummyAccounts()
	return &s
}

func (suite *KeeperTestSuite) initDummyAccounts() {
	_ = suite.MintCoins(ALICE, 1000*KYVE)
	_ = suite.MintCoins(BOB, 1000*KYVE)
	_ = suite.MintCoins(CHARLIE, 1000*KYVE)
	_ = suite.MintCoins(DAVID, 1000*KYVE)

	_ = suite.MintCoins(STAKER_0, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_0_A, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_0_B, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_0_C, 1000*KYVE)

	_ = suite.MintCoins(STAKER_1, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_1_A, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_1_B, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_1_C, 1000*KYVE)

	_ = suite.MintCoins(STAKER_2, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_2_A, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_2_B, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_2_C, 1000*KYVE)

	_ = suite.MintCoins(STAKER_3, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_3_A, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_3_B, 1000*KYVE)
	_ = suite.MintCoins(VALADDRESS_3_C, 1000*KYVE)

	DUMMY = make([]string, 50)

	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			randomSource := mrand.New(mrand.NewSource(int64(i + k)))
			byteAddr[k] = byte(randomSource.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("kyve", byteAddr)
		DUMMY[i] = dummy
		_ = suite.MintCoins(dummy, 1000*KYVE)
	}

	VALDUMMY = make([]string, 50)
	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			randomSource := mrand.New(mrand.NewSource(int64(i + k + 100)))
			byteAddr[k] = byte(randomSource.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("kyve", byteAddr)
		VALDUMMY[i] = dummy
		_ = suite.MintCoins(dummy, 1000*KYVE)
	}
}

func (suite *KeeperTestSuite) MintCoins(address string, amount uint64) error {
	// mint coins ukyve, A, B, C
	coins := sdk.NewCoins(
		sdk.NewInt64Coin(KYVE_DENOM, int64(amount)),
		sdk.NewInt64Coin(A_DENOM, int64(amount)),
		sdk.NewInt64Coin(B_DENOM, int64(amount)),
		sdk.NewInt64Coin(C_DENOM, int64(amount)),
	)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, mintTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	receiver, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, mintTypes.ModuleName, receiver, coins)
	if err != nil {
		return err
	}

	return nil
}

func (suite *KeeperTestSuite) MintCoin(address string, coin sdk.Coin) error {
	// mint coins ukyve, A, B, C
	coins := sdk.NewCoins(coin)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, mintTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	receiver, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, mintTypes.ModuleName, receiver, coins)
	if err != nil {
		return err
	}

	return nil
}

func (suite *KeeperTestSuite) MintDenomToModule(moduleAddress string, amount uint64, denom string) error {
	coins := sdk.NewCoins(sdk.NewInt64Coin(denom, int64(amount)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, mintTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, mintTypes.ModuleName, moduleAddress, coins)
	if err != nil {
		return err
	}

	return nil
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app         *app.App
	address     [20]byte
	consAddress sdk.ConsAddress
	validator   stakingtypes.Validator
	denom       string
}

func (suite *KeeperTestSuite) App() *app.App {
	return suite.app
}

func (suite *KeeperTestSuite) Ctx() sdk.Context {
	return suite.ctx
}

func (suite *KeeperTestSuite) SetCtx(ctx sdk.Context) {
	suite.ctx = ctx
}

func (suite *KeeperTestSuite) SetupTest(startTime int64) {
	suite.SetupApp(startTime)
}

func (suite *KeeperTestSuite) SetupApp(startTime int64) {
	suite.app = app.Setup()

	suite.denom = globalTypes.Denom

	rawHex, _ := hex.DecodeString("0xBf71F763e4DEd30139C40160AE74Df881D5C7A2d")
	suite.address = [20]byte(rawHex[:20])

	// consensus key
	ePriv := ed25519.GenPrivKeyFromSecret([]byte{1})
	suite.consAddress = sdk.ConsAddress(ePriv.PubKey().Address())

	suite.ctx = suite.app.BaseApp.NewContextLegacy(false, tmproto.Header{
		Height:          1,
		ChainID:         "kyve-test",
		Time:            time.Unix(startTime, 0).UTC(),
		ProposerAddress: suite.consAddress.Bytes(),

		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	})

	mintParams, _ := suite.app.MintKeeper.Params.Get(suite.ctx)
	mintParams.MintDenom = suite.denom
	_ = suite.app.MintKeeper.Params.Set(suite.ctx, mintParams)

	stakingParams, _ := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = suite.denom
	_ = suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	govParams, _ := suite.app.GovKeeper.Params.Get(suite.ctx)
	govParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(KYVE_DENOM, int64(100_000_000_000))) // set min deposit to 100 KYVE
	_ = suite.app.GovKeeper.Params.Set(suite.ctx, govParams)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address[:])
	validator, _ := stakingtypes.NewValidator(valAddr.String(), ePriv.PubKey(), stakingtypes.Description{})
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	//_ = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	_ = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	validators, _ := suite.app.StakingKeeper.GetValidators(suite.ctx, 1)
	suite.validator = validators[0]
}

func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

func (suite *KeeperTestSuite) CommitAfterSeconds(seconds uint64) {
	suite.CommitAfter(time.Second * time.Duration(seconds))
}

func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
	header := suite.ctx.BlockHeader()
	header.Time = header.Time.Add(t)

	_, err := suite.app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: header.Height,
		Time:   header.Time,
	})
	if err != nil {
		panic(err)
	}
	_, err = suite.app.Commit()
	if err != nil {
		panic(err)
	}

	header.Height += 1

	suite.ctx = suite.app.BaseApp.NewUncachedContext(false, header)
}

func (suite *KeeperTestSuite) WaitSeconds(seconds uint64) {
	suite.Wait(time.Second * time.Duration(seconds))
}

func (suite *KeeperTestSuite) Wait(t time.Duration) {
	header := suite.ctx.BlockHeader()
	header.Time = header.Time.Add(t)

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(t))
}
