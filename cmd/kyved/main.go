package main

import (
	"os"

	kyveApp "github.com/KYVENetwork/chain/app"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	initSDKConfig("kyve")
	rootCmd := NewRootCmd(kyveApp.MakeEncodingConfig())
	if err := serverCmd.Execute(rootCmd, "", kyveApp.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
