package cli

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdToggleMultiCoinRewards())
	cmd.AddCommand(CmdSetMultiCoinDistributionPolicy())

	return cmd
}
