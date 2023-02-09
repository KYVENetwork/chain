package app

import (
	"encoding/json"
	"time"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	teamTypes "github.com/KYVENetwork/chain/x/team/types"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	tmtypes "github.com/tendermint/tendermint/types"
)

// DefaultConsensusParams ...
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   -1, // no limit
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func DefaultGenesisWithValSet(codec codec.Codec) map[string]json.RawMessage {
	bondingDenom := globalTypes.Denom

	// Generate a new validator.
	key, _ := mock.NewPV().GetPubKey()
	validator := tmtypes.NewValidator(key, 1)

	publicKey, _ := cryptocodec.FromTmPubKeyInterface(validator.PubKey)
	publicKeyAny, _ := codectypes.NewAnyWithValue(publicKey)

	validators := []stakingtypes.Validator{
		{
			OperatorAddress:   sdk.ValAddress(validator.Address).String(),
			ConsensusPubkey:   publicKeyAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            sdk.DefaultPowerReduction,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   0,
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		},
	}
	// Generate a new delegator.
	delegatorKey := secp256k1.GenPrivKey()
	delegator := authtypes.NewBaseAccount(
		delegatorKey.PubKey().Address().Bytes(), delegatorKey.PubKey(), 0, 0,
	)

	delegations := []stakingtypes.Delegation{
		stakingtypes.NewDelegation(delegator.GetAddress(), validator.Address.Bytes(), sdk.OneDec()),
	}

	// Default genesis state.
	config := MakeEncodingConfig()
	genesisState := ModuleBasics.DefaultGenesis(config.Marshaler)

	// Update x/auth state.
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), []authtypes.GenesisAccount{delegator})
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	// Update x/bank state.
	bondedCoins := sdk.NewCoins(sdk.NewCoin(bondingDenom, sdk.DefaultPowerReduction))

	teamCoins := sdk.NewCoins(sdk.NewInt64Coin(globalTypes.Denom, int64(teamTypes.TEAM_ALLOCATION)))

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, []banktypes.Balance{
		{
			Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
			Coins:   bondedCoins,
		},
		{
			Address: authtypes.NewModuleAddress(teamTypes.ModuleName).String(),
			Coins:   teamCoins,
		},
	}, bondedCoins.Add(sdk.NewInt64Coin(globalTypes.Denom, int64(teamTypes.TEAM_ALLOCATION))), []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	// Update x/staking state.
	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = bondingDenom

	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	// Return.
	return genesisState
}

// Setup initializes a new App.
func Setup() *App {
	db := dbm.NewMemDB()

	config := MakeEncodingConfig()

	setPrefixes("kyve")

	app := NewKYVEApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, config, simapp.EmptyAppOptions{})
	// init chain must be called to stop deliverState from being nil

	genesisState := DefaultGenesisWithValSet(app.AppCodec())
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         "kyve-test",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

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
}
