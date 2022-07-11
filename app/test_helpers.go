package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/ignite-hq/cli/ignite/pkg/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"time"

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

// Setup initializes a new App.
func Setup() *App {
	db := dbm.NewMemDB()

	config := cosmoscmd.MakeEncodingConfig(ModuleBasics)

	app := NewApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, config, simapp.EmptyAppOptions{})
	// init chain must be called to stop deliverState from being nil
	genesisState := ModuleBasics.DefaultGenesis(config.Marshaler)

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
