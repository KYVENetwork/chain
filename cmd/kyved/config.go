package main

import (
	"fmt"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"
	serverCfg "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmCfg "github.com/tendermint/tendermint/config"
)

func initAppConfig() (string, *serverCfg.Config) {
	cfg := serverCfg.DefaultConfig()
	cfg.MinGasPrices = fmt.Sprintf("0.001%s", globalTypes.Denom)

	return serverCfg.DefaultConfigTemplate, cfg
}

func initSDKConfig(accountAddressPrefix string) {
	accountPubKeyPrefix := accountAddressPrefix + "pub"
	validatorAddressPrefix := accountAddressPrefix + "valoper"
	validatorPubKeyPrefix := accountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := accountAddressPrefix + "valcons"
	consNodePubKeyPrefix := accountAddressPrefix + "valconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(accountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}

func initTendermintConfig() *tmCfg.Config {
	return tmCfg.DefaultConfig()
}
