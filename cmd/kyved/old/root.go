package old

import (
	kyveApp "github.com/KYVENetwork/chain/app"
	tmCli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"os"

	// Auth
	authCli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// Bank
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	// Crisis
	"github.com/cosmos/cosmos-sdk/x/crisis"
	// GenUtil
	genUtilCli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genUtilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// NewRootCmd creates a new root command for the KYVE chain daemon.
func NewRootCmd(encodingConfig kyveApp.EncodingConfig) *cobra.Command {
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithHomeDir(kyveApp.DefaultNodeHome).
		WithViper("KYVE")

	rootCmd := &cobra.Command{
		Use:   "kyved",
		Short: "KYVE Chain Daemon",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := initTendermintConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
	}

	ac := appCreator{encodingConfig}
	server.AddCommands(
		rootCmd,
		kyveApp.DefaultNodeHome,
		ac.createApp,
		ac.exportApp,
		func(startCmd *cobra.Command) {
			crisis.AddModuleInitFlags(startCmd)
		},
	)

	rootCmd.AddCommand(
		genUtilCli.InitCmd(kyveApp.ModuleBasics, kyveApp.DefaultNodeHome),
		// TODO(@john): Investigate why the one directly from the module is nil.
		genUtilCli.CollectGenTxsCmd(bankTypes.GenesisBalancesIterator{}, kyveApp.DefaultNodeHome, genUtilTypes.DefaultMessageValidator, ac.encodingConfig.AddressCoded),
		// TODO(@rapha): fix migrations
		genUtilCli.MigrateGenesisCmd(nil),
		genUtilCli.GenTxCmd(
			kyveApp.ModuleBasics,
			encodingConfig.TxConfig,
			bankTypes.GenesisBalancesIterator{},
			kyveApp.DefaultNodeHome,
			ac.encodingConfig.AddressCoded,
		),
		infoCommand(),
		genUtilCli.ValidateGenesisCmd(kyveApp.ModuleBasics),
		addGenesisAccountCmd(kyveApp.DefaultNodeHome),
		tmCli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		// TODO(@rapha): fix StatusCommand
		//config.Cmd(),
		pruning.Cmd(ac.createApp, kyveApp.DefaultNodeHome),

		// TODO(@rapha): fix StatusCommand
		//rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		// TODO(@rapha): fix this
		//keys.Commands(kyveApp.DefaultNodeHome),
	)

	return rootCmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		// TODO(@rapha): fix this
		//rpc.BlockCommand(),
		rpc.ValidatorCommand(),
		// TODO(@rapha): fix this
		//authCli.GetAccountCmd(),
		authCli.QueryTxCmd(),
		authCli.QueryTxsByEventsCmd(),
	)

	kyveApp.ModuleBasics.AddQueryCommands(cmd)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authCli.GetSignCommand(),
		authCli.GetSignBatchCommand(),
		authCli.GetMultiSignCommand(),
		authCli.GetValidateSignaturesCommand(),
		authCli.GetBroadcastCommand(),
		authCli.GetEncodeCommand(),
		authCli.GetDecodeCommand(),
	)

	kyveApp.ModuleBasics.AddTxCommands(cmd)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}
