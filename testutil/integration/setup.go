package integration

import (
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"encoding/json"
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
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"time"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	app   *app.App
	denom string
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

	setPrefixes(app.AccountAddressPrefix)

	logger := log.NewNopLogger()
	localApp, err := app.New(logger, db, nil, true, EmptyAppOptions{}, baseapp.SetChainID("kyve-test"))
	if err != nil {
		panic(err)
	}
	suite.app = localApp

	genesisState := DefaultGenesisWithValSet(suite.app)
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
		ProposerAddress: sdk.ConsAddress(ed25519.GenPrivKeyFromSecret([]byte("Validator-1")).PubKey().Address()).Bytes(),

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
	stakingParams.MaxValidators = 51
	_ = suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)

	govParams, _ := suite.app.GovKeeper.Params.Get(suite.ctx)
	govParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(KYVE_DENOM, int64(100_000_000_000))) // set min deposit to 100 KYVE
	_ = suite.app.GovKeeper.Params.Set(suite.ctx, govParams)
}

func setPrefixes(accountAddressPrefix string) {
	// Set prefixes
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	config := sdk.GetConfig()

	// Return if prefixes are already set
	if config.GetBech32AccountAddrPrefix() == accountAddressPrefix &&
		config.GetBech32AccountPubPrefix() == accountPubKeyPrefix {
		return
	}

	// Set and seal config
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}

func DefaultGenesisWithValSet(app *app.App) map[string]json.RawMessage {
	bondingDenom := globalTypes.Denom

	// Generate a new validator.
	pubKey := ed25519.GenPrivKey().PubKey()
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
	delegatorKey := secp256k1.GenPrivKey()
	delegator := authTypes.NewBaseAccount(
		delegatorKey.PubKey().Address().Bytes(), delegatorKey.PubKey(), 0, 0,
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

	stakingGenesis := stakingTypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingTypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	// Return.
	return genesisState
}

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} { return nil }
