package main

import (
	"errors"
	"io"

	kyveApp "github.com/KYVENetwork/chain/app"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/viper"
)

// appCreator is a wrapper for EncodingConfig.
// This allows us to reuse encodingConfig received by NewRootCmd in both createApp and exportApp.
type appCreator struct{ encodingConfig kyveApp.EncodingConfig }

func (ac appCreator) createApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts serverTypes.AppOptions,
) serverTypes.Application {
	return kyveApp.NewKYVEApp(
		logger, db, traceStore, true,
		appOpts,
		server.DefaultBaseappOptions(appOpts)...,
	)
}

func (ac appCreator) exportApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts serverTypes.AppOptions,
	modulesToExport []string,
) (serverTypes.ExportedApp, error) {
	var app *kyveApp.KYVEApp

	// this check is necessary as we use the flag in x/upgrade.
	// we can exit more gracefully by checking the flag here.
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return serverTypes.ExportedApp{}, errors.New("application home not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return serverTypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	if height != -1 {
		app = kyveApp.NewKYVEApp(logger, db, traceStore, false, appOpts)

		if err := app.LoadHeight(height); err != nil {
			return serverTypes.ExportedApp{}, err
		}
	} else {
		app = kyveApp.NewKYVEApp(logger, db, traceStore, true, appOpts)
	}

	return app.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
