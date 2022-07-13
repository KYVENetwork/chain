package keeper_test

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/KYVENetwork/chain/app"
	"github.com/KYVENetwork/chain/x/registry"
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	mrand "math/rand"
	"testing"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/version"
	"time"
)

const ALICE_ADDR = "cosmos1jq304cthpx0lwhpqzrdjrcza559ukyy347ju8f"
const BOB_ADDR = "cosmos1hvg7zsnrj6h29q9ss577mhrxa04rn94hfvl2ry"

//const ALICE_ADDR = "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"
//const BOB_ADDR = "kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq"

const KYVE = uint64(1_000_000_000)

var DUMMY_ACCOUNTS []string

func runTxWithResult(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := s.ctx.CacheContext()
	resp, err := registry.NewHandler(s.app.RegistryKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func runTx(msg sdk.Msg) (success bool) {
	cachedCtx, commit := s.ctx.CacheContext()
	_, err := registry.NewHandler(s.app.RegistryKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return true
	}
	return false
}

func runTxSuccess(t *testing.T, msg sdk.Msg) {
	success := runTx(msg)
	require.True(t, success)
}

func mint(address string, amount uint64) error {
	coins := sdk.NewCoins(sdk.NewInt64Coin("tkyve", int64(amount)))
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}

	s.Commit()

	//accPrefix := cosmostypes.GetConfig().GetBech32AccountAddrPrefix()
	//pubPrefix := cosmostypes.GetConfig().GetBech32AccountPubPrefix()
	//cosmostypes.GetConfig().SetBech32PrefixForAccount("kyve", pubPrefix)
	sender, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}
	//cosmostypes.GetConfig().SetBech32PrefixForAccount(accPrefix, pubPrefix)

	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, sender, coins)
	if err != nil {
		return err
	}
	return nil
}

func initDummyAccounts() {
	DUMMY_ACCOUNTS = make([]string, 50)
	mrand.Seed(1)
	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			byteAddr[k] = byte(mrand.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("cosmos", byteAddr)
		DUMMY_ACCOUNTS[i] = dummy
		mint(dummy, 1000*KYVE)
	}
}

func createGenesis(t *testing.T) {
	s = new(KeeperTestSuite)
	s.SetupTest()

	currentTime := s.ctx.BlockTime().Unix()
	s.CommitAfter(time.Second * 60)
	require.Equal(t, s.ctx.BlockTime().Unix(), currentTime+60)

	s.CommitAfter(time.Second * 60)
	require.Equal(t, s.ctx.BlockTime().Unix(), currentTime+2*60)

	pool := types.Pool{
		Creator:        govtypes.ModuleName,
		Name:           "Moontest",
		Runtime:        "@kyve/evm",
		Logo:           "9FJDam56yBbmvn8rlamEucATH5UcYqSBw468rlCXn8E",
		Config:         "{\"rpc\":\"https://rpc.api.moonbeam.network\",\"github\":\"https://github.com/KYVENetwork/evm\"}",
		UploadInterval: 60,
		OperatingCost:  100,
		BundleProposal: &types.BundleProposal{},
		MaxBundleSize:  100,
		Protocol: &types.Protocol{
			Version:     "1.3.0",
			LastUpgrade: uint64(s.ctx.BlockTime().Unix()),
			Binaries:    "{\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.0.5/kyve-evm-macos.zip\"}",
		},
		UpgradePlan: &types.UpgradePlan{},
		StartKey:    "0",
		Status:      types.POOL_STATUS_NOT_ENOUGH_VALIDATORS,
		MinStake:    0,
	}

	s.app.RegistryKeeper.AppendPool(s.ctx, pool)
	s.Commit()

	err := mint(ALICE_ADDR, 1000*KYVE)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	err = mint(BOB_ADDR, 1000*KYVE)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	initDummyAccounts()
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app         *app.App
	queryClient types.QueryClient
	address     common.Address
	signer      keyring.Signer
	consAddress sdk.ConsAddress
	validator   stakingtypes.Validator
	denom       string
}

var s *KeeperTestSuite

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = app.Setup()
	suite.SetupApp()
}

func (suite *KeeperTestSuite) SetupApp() {
	//t := suite.T()

	suite.denom = "tkyve"

	//fmt.Printf("%s\n", sdk.GetConfig().GetBech32AccountAddrPrefix())
	//fmt.Printf("%s\n", sdk.GetConfig().GetBech32AccountPubPrefix())
	//fmt.Printf("%s\n", sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	//fmt.Printf("%s\n", sdk.GetConfig().GetBech32ValidatorPubPrefix())
	//fmt.Printf("%s\n", sdk.GetConfig().GetBech32ConsensusAddrPrefix())
	//fmt.Printf("%s\n", sdk.GetConfig().GetBech32ConsensusPubPrefix())

	//sdk.GetConfig().SetBech32PrefixForAccount("kyve", "kyvepub")
	//sdk.GetConfig().SetBech32PrefixForValidator("kyvevaloper", "kyvevaloperpub")
	//sdk.GetConfig().SetBech32PrefixForValidator("kyvevalcons", "kyvevalconspub")

	// consensus key
	privKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	_ = err
	//require.NoError(t, err)

	addressBytes := tmcrypto.Address(crypto.PubkeyToAddress(privKey.PublicKey).Bytes())
	suite.address = common.BytesToAddress(addressBytes)

	// consensus key
	privKey, err = ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	//require.NoError(t, err)

	ePriv := ed25519.GenPrivKeyFromSecret([]byte{1})
	suite.consAddress = sdk.ConsAddress(ePriv.PubKey().Address())

	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Height:          1,
		ChainID:         "kyve-test",
		Time:            time.Now().UTC(),
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
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.RegistryKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	mintParams := suite.app.MintKeeper.GetParams(suite.ctx)
	mintParams.MintDenom = suite.denom
	suite.app.MintKeeper.SetParams(suite.ctx, mintParams)

	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = suite.denom
	suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	// Set Validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, ePriv.PubKey(), stakingtypes.Description{})
	//require.NoError(t, err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper, suite.ctx, validator, true)
	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	//require.NoError(t, err)
	validators := s.app.StakingKeeper.GetValidators(s.ctx, 1)
	suite.validator = validators[0]
}

func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

func (suite *KeeperTestSuite) CommitAfterSeconds(seconds uint64) {
	suite.CommitAfter(time.Second * time.Duration(seconds))
}

// Commit commits a block at a given time.
func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
	header := suite.ctx.BlockHeader()
	suite.app.EndBlock(abci.RequestEndBlock{Height: header.Height})
	_ = suite.app.Commit()

	header.Height += 1
	header.Time = header.Time.Add(t)
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// update ctx
	suite.ctx = suite.app.BaseApp.NewContext(false, header)

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.RegistryKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}
