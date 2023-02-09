package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KYVENetwork/chain/x/stakers/types"
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

	cmd.AddCommand(CmdCreateStaker())
	cmd.AddCommand(CmdJoinPool())
	cmd.AddCommand(CmdLeavePool())
	cmd.AddCommand(CmdUpdateCommission())
	cmd.AddCommand(CmdUpdateMetadata())

	return cmd
}
