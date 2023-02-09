package cli

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/team/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
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

	cmd.AddCommand(CmdClaimUnlocked())
	cmd.AddCommand(CmdClawback())
	cmd.AddCommand(CmdCreateTeamVestingAccount())
	cmd.AddCommand(CmdClaimAuthorityRewards())
	cmd.AddCommand(CmdClaimAccountRewards())

	return cmd
}
