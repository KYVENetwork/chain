package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KYVENetwork/chain/x/query/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group query queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// Account
	cmd.AddCommand(CmdAccountAssets())
	cmd.AddCommand(CmdAccountFundedList())
	cmd.AddCommand(CmdAccountDelegationUnbondings())
	cmd.AddCommand(CmdAccountRedelegation())

	// Pool
	cmd.AddCommand(CmdShowPool())
	cmd.AddCommand(CmdListPool())

	// Staking
	cmd.AddCommand(CmdShowStaker())
	cmd.AddCommand(CmdListStakers())
	cmd.AddCommand(CmdListStakersByPool())

	// DELEGATION
	cmd.AddCommand(CmdDelegator())
	cmd.AddCommand(CmdStakersByPoolAndDelegator())
	cmd.AddCommand(CmdDelegatorsByPoolAndStaker())

	// Bundles
	cmd.AddCommand(CmdShowFinalizedBundle())
	cmd.AddCommand(CmdListFinalizedBundles())
	cmd.AddCommand(CmdCanPropose())
	cmd.AddCommand(CmdCanVote())
	cmd.AddCommand(CmdCurrentVoteStatus())
	cmd.AddCommand(CmdCanValidate())

	return cmd
}
