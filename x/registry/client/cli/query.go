package cli

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group registry queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// PARAMS
	cmd.AddCommand(CmdParams())

	// POOL
	cmd.AddCommand(CmdShowPool())
	cmd.AddCommand(CmdListPool())
	cmd.AddCommand(CmdFundersList())
	cmd.AddCommand(CmdFunder())
	cmd.AddCommand(CmdStakersList())
	cmd.AddCommand(CmdStaker())

	// WARP
	cmd.AddCommand(CmdShowProposal())
	cmd.AddCommand(CmdListProposal())
	cmd.AddCommand(CmdProposalByHeight())

	// PROTOCOL NODE - FLOW
	cmd.AddCommand(CmdCanPropose())
	cmd.AddCommand(CmdCanVote())
	cmd.AddCommand(CmdStakeInfo())

	// STATS FOR USER ACCOUNT
	cmd.AddCommand(CmdAccountAssets())
	cmd.AddCommand(CmdAccountFundedList())
	cmd.AddCommand(CmdAccountStakedList())
	cmd.AddCommand(CmdAccountDelegationList())

	// DELEGATION
	cmd.AddCommand(CmdDelegator())
	cmd.AddCommand(CmdStakersByPoolAndDelegator())
	cmd.AddCommand(CmdDelegatorsByPoolAndStaker())

	// this line is used by starport scaffolding # 1

	return cmd
}
