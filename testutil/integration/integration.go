package integration

import (
	mrand "math/rand"
	"time"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/app"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

const (
	ALICE   = "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"
	BOB     = "kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq"
	CHARLIE = "kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s"

	STAKER_0     = "kyve1htgfatqevuvfzvl0sxp97ywteqhg5leha9emf4"
	VALADDRESS_0 = "kyve1qnf86dkvvtpdukx30r3vajav7rdq8snktm90hm"

	STAKER_1     = "kyve1gnr35rwn8rmflnlzs6nn5hhkmzzkxg9ap8xepw"
	VALADDRESS_1 = "kyve1hpjgzljglmv00nstk3jvcw0zzq94nu0cuxv5ga"

	STAKER_2     = "kyve1xsemlxghgvusumhqzm2ztjw7dz9krvu3de54e2"
	VALADDRESS_2 = "kyve1u0870dkae6ql63hxvy9y7g65c0y8csfh8allzl"

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

	_ = suite.Mint(STAKER_0, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_0, 1000*KYVE)

	_ = suite.Mint(STAKER_1, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_1, 1000*KYVE)

	_ = suite.Mint(STAKER_2, 1000*KYVE)
	_ = suite.Mint(VALADDRESS_2, 1000*KYVE)

	mrand.New(mrand.NewSource(1))

	DUMMY = make([]string, 50)

	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			mrand.New(mrand.NewSource(int64(i + k)))
			byteAddr[k] = byte(mrand.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("kyve", byteAddr)
		DUMMY[i] = dummy
		_ = suite.Mint(dummy, 1000*KYVE)
	}

	VALDUMMY = make([]string, 50)
	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			mrand.New(mrand.NewSource(int64(i + k + 100)))
			byteAddr[k] = byte(mrand.Int())
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
	suite.app.MintKeeper.SetParams(suite.ctx, mintParams)

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = suite.denom
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	depositParams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
	depositParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(KYVE_DENOM, int64(100_000_000_000))) // set min deposit to 100 KYVE
	suite.app.GovKeeper.SetDepositParams(suite.ctx, depositParams)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, _ := stakingtypes.NewValidator(valAddr, ePriv.PubKey(), stakingtypes.Description{})
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	_ = suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
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
