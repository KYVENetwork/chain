package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types/module"
	"time"

	"cosmossdk.io/math"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtProto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmtTypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// Auth
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Global
	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	// Staking
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// Team
	teamTypes "github.com/KYVENetwork/chain/x/team/types"
)

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

type EmptyAppOptions struct{}

func (ao EmptyAppOptions) Get(_ string) interface{} { return nil }

func DefaultGenesisWithValSet(app *App) map[string]json.RawMessage {
	bondingDenom := globalTypes.Denom

	// Generate a new validator.
	key, _ := mock.NewPV().GetPubKey()
	validator := cmtTypes.NewValidator(key, 1)

	publicKey, _ := cryptoCodec.FromTmPubKeyInterface(validator.PubKey)
	publicKeyAny, _ := codecTypes.NewAnyWithValue(publicKey)

	validators := []stakingTypes.Validator{
		{
			OperatorAddress:   sdk.ValAddress(validator.Address).String(),
			ConsensusPubkey:   publicKeyAny,
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
		stakingTypes.NewDelegation(delegator.GetAddress().String(), validator.Address.String(), math.LegacyOneDec()),
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

// Setup initializes a new App.
func Setup() *App {
	db := dbm.NewMemDB()

	// config := MakeEncodingConfig()

	setPrefixes("kyve")

	// app := NewKYVEApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, config, EmptyAppOptions{})
	app, err := New(log.NewNopLogger(), db, nil, true, EmptyAppOptions{}, baseapp.SetChainID("kyve-test"))
	if err != nil {
		panic(err)
	}

	// TODO: Do we need this?
	kyveModules := RegisterKyveModules(app.InterfaceRegistry())
	for name, mod := range kyveModules {
		app.ModuleManager.Modules[name] = module.CoreAppModuleBasicAdaptor(name, mod)
		//app.autoCliOpts.Modules[name] = mod
	}

	genesisState := DefaultGenesisWithValSet(app)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	_, err = app.InitChain(
		&abci.RequestInitChain{
			ChainId:         "kyve-test",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	if err != nil {
		panic(err)
	}

	return app
}

func setPrefixes(accountAddressPrefix string) {
	// Set prefixes
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}
