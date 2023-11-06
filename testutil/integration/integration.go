package integration

import (
	mrand "math/rand"
	"time"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/app"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	"github.com/cometbft/cometbft/version"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
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
	KYVE  = uint64(1_000_000_000)
	TKYVE = uint64(1)
)

var KYVE_DENOM = globalTypes.Denom

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
	_ = suite.Mint(ALICE, 1000*KYVE)
	_ = suite.Mint(BOB, 1000*KYVE)
	_ = suite.Mint(CHARLIE, 1000*KYVE)
	_ = suite.Mint(DAVID, 1000*KYVE)

	_ = suite.Mint(STAKER_0, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_0_A, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_0_B, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_0_C, 1000*KYVE)

	_ = suite.Mint(STAKER_1, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_1_A, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_1_B, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_1_C, 1000*KYVE)

	_ = suite.Mint(STAKER_2, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_2_A, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_2_B, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_2_C, 1000*KYVE)

	_ = suite.Mint(STAKER_3, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_3_A, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_3_B, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_3_C, 1000*KYVE)

	DUMMY = make([]string, 50)

	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			randomSource := mrand.New(mrand.NewSource(int64(i + k)))
			byteAddr[k] = byte(randomSource.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("kyve", byteAddr)
		DUMMY[i] = dummy
		_ = suite.Mint(dummy, 1000*KYVE)
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
		_ = suite.Mint(dummy, 1000*KYVE)
	}
}

func (suite *KeeperTestSuite) Mint(address string, amount uint64) error {
	coins := sdk.NewCoins(sdk.NewInt64Coin(KYVE_DENOM, int64(amount)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, mintTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	sender, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, mintTypes.ModuleName, sender, coins)
	if err != nil {
		return err
	}

	return nil
}

type QueryClients struct {
	stakersClient stakerstypes.QueryClient
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app         *app.App
	queries     QueryClients
	address     common.Address
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

	suite.address = common.HexToAddress("0xBf71F763e4DEd30139C40160AE74Df881D5C7A2d")

	// consensus key
	ePriv := ed25519.GenPrivKeyFromSecret([]byte{1})
	suite.consAddress = sdk.ConsAddress(ePriv.PubKey().Address())

	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
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
	suite.registerQueryClients()

	mintParams := suite.app.MintKeeper.GetParams(suite.ctx)
	mintParams.MintDenom = suite.denom
	_ = suite.app.MintKeeper.SetParams(suite.ctx, mintParams)

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = suite.denom
	_ = suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	govParams := suite.app.GovKeeper.GetParams(suite.ctx)
	govParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(KYVE_DENOM, int64(100_000_000_000))) // set min deposit to 100 KYVE
	_ = suite.app.GovKeeper.SetParams(suite.ctx, govParams)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, _ := stakingtypes.NewValidator(valAddr, ePriv.PubKey(), stakingtypes.Description{})
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	//_ = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	_ = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	validators := suite.app.StakingKeeper.GetValidators(suite.ctx, 1)
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
	suite.app.EndBlock(abci.RequestEndBlock{Height: header.Height})
	_ = suite.app.Commit()

	header.Height += 1
	header.Time = header.Time.Add(t)
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: header})

	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	suite.registerQueryClients()
}

func (suite *KeeperTestSuite) registerQueryClients() {
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())

	stakerstypes.RegisterQueryServer(queryHelper, suite.app.StakersKeeper)
	suite.queries.stakersClient = stakerstypes.NewQueryClient(queryHelper)
}
