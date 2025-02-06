package integration

import (
	"encoding/json"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/KYVENetwork/chain/app"
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	cmtProto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cmtTypes "github.com/cometbft/cometbft/types"
	"github.com/cometbft/cometbft/version"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app                 *app.App
	denom               string
	privateValidatorKey *ed25519.PrivKey
	VoteInfos           []abci.VoteInfo
}

// DefaultConsensusParams ...
var DefaultConsensusParams = &cmtProto.ConsensusParams{
	Block: &cmtProto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1, // no limit
	},
	Evidence: &cmtProto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &cmtProto.ValidatorParams{
		PubKeyTypes: []string{
			cmtTypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func (suite *KeeperTestSuite) SetupApp(startTime int64) {
	db := dbm.NewMemDB()

	config := sdk.GetConfig()
	if config.GetBech32AccountAddrPrefix() != "kyve" ||
		config.GetBech32AccountPubPrefix() != "kyvepub" {
		config.SetBech32PrefixForAccount("kyve", "kyve"+"pub")
		config.SetBech32PrefixForValidator("kyve"+"valoper", "kyve"+"valoperpub")
		config.SetBech32PrefixForConsensusNode("kyve"+"valcons", "kyve"+"valconspub")
		config.Seal()
	}

	logger := log.NewNopLogger()
	localApp, err := app.New(logger, db, nil, true, EmptyAppOptions{}, baseapp.SetChainID("kyve-test"))
	if err != nil {
		panic(err)
	}
	suite.app = localApp

	suite.privateValidatorKey = ed25519.GenPrivKeyFromSecret([]byte("Validator-1"))
	genesisState := DefaultGenesisWithValSet(suite.app, suite.privateValidatorKey)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	if _, err = suite.app.InitChain(
		&abci.RequestInitChain{
			ChainId:         "kyve-test",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	); err != nil {
		panic(err)
	}

	suite.denom = globalTypes.Denom

	suite.ctx = suite.app.BaseApp.NewContextLegacy(false, tmproto.Header{
		Height:          1,
		ChainID:         "kyve-test",
		Time:            time.Unix(startTime, 0).UTC(),
		ProposerAddress: sdk.ConsAddress(suite.privateValidatorKey.PubKey().Address()).Bytes(),

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
}

func DefaultGenesisWithValSet(app *app.App, validatorPrivateKey *ed25519.PrivKey) map[string]json.RawMessage {
	bondingDenom := globalTypes.Denom

	// Generate a new validator.
	pubKey := validatorPrivateKey.PubKey()
	valAddress := sdk.ValAddress(pubKey.Address()).String()
	pkAny, _ := codectypes.NewAnyWithValue(pubKey)

	validators := []stakingTypes.Validator{
		{
			OperatorAddress:   valAddress,
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingTypes.Bonded,
			Tokens:            sdk.DefaultPowerReduction,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingTypes.Description{},
			UnbondingHeight:   0,
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingTypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		},
	}
	// Generate a new delegator.
	delegator := authTypes.NewBaseAccount(
		validatorPrivateKey.PubKey().Address().Bytes(), validatorPrivateKey.PubKey(), 0, 0,
	)

	delegations := []stakingTypes.Delegation{
		stakingTypes.NewDelegation(delegator.GetAddress().String(), valAddress, math.LegacyOneDec()),
	}

	// Default genesis state.
	genesisState := app.DefaultGenesis()

	// Update x/auth state.
	authGenesis := authTypes.NewGenesisState(authTypes.DefaultParams(), []authTypes.GenesisAccount{delegator})
	genesisState[authTypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	// Update x/bank state.
	bondedCoins := sdk.NewCoins(sdk.NewCoin(bondingDenom, sdk.DefaultPowerReduction))

	teamCoins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(teamTypes.TEAM_ALLOCATION)))

	bankGenesis := bankTypes.NewGenesisState(bankTypes.DefaultGenesisState().Params, []bankTypes.Balance{
		{
			Address: authTypes.NewModuleAddress(stakingTypes.BondedPoolName).String(),
			Coins:   bondedCoins,
		},
		{
			Address: authTypes.NewModuleAddress(teamTypes.ModuleName).String(),
			Coins:   teamCoins,
		},
	}, bondedCoins.Add(sdk.NewInt64Coin(globalTypes.Denom, int64(teamTypes.TEAM_ALLOCATION))), []bankTypes.Metadata{}, []bankTypes.SendEnabled{})
	genesisState[bankTypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Update x/staking state.
	stakingParams := stakingTypes.DefaultParams()
	stakingParams.BondDenom = bondingDenom
	stakingParams.MaxValidators = 51

	stakingGenesis := stakingTypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingTypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	mintParams := mintTypes.DefaultParams()
	mintParams.MintDenom = bondingDenom
	mintGenesis := mintTypes.NewGenesisState(mintTypes.DefaultInitialMinter(), mintParams)
	genesisState[mintTypes.ModuleName] = app.AppCodec().MustMarshalJSON(mintGenesis)

	govParams := govTypes.DefaultParams()
	govParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(KYVE_DENOM, int64(100*KYVE)))
	govGenesis := govTypes.NewGenesisState(1, govParams)
	genesisState["gov"] = app.AppCodec().MustMarshalJSON(govGenesis)

	slashingParams := slashingTypes.DefaultParams()
	slashingParams.SignedBlocksWindow = 1_000_000
	// Allow always being offline
	slashingParams.MinSignedPerWindow = math.LegacyMustNewDecFromStr("0.1")
	slashingInfo := []slashingTypes.SigningInfo{{
		Address: sdk.MustBech32ifyAddressBytes("kyvevalcons", validatorPrivateKey.PubKey().Address()),
		ValidatorSigningInfo: slashingTypes.ValidatorSigningInfo{
			Address:             sdk.MustBech32ifyAddressBytes("kyvevalcons", validatorPrivateKey.PubKey().Address()),
			StartHeight:         0,
			IndexOffset:         0,
			JailedUntil:         time.Time{},
			Tombstoned:          false,
			MissedBlocksCounter: 0,
		},
	}}
	slashingGenesis := slashingTypes.NewGenesisState(slashingParams, slashingInfo, nil)
	genesisState["slashing"] = app.AppCodec().MustMarshalJSON(slashingGenesis)

	return genesisState
}

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} { return nil }
