package cli

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdVoteBundleProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-bundle-proposal [staker] [pool_id] [storage_id] [vote]",
		Short: "Broadcast message vote-bundle-proposal",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStaker := args[0]

			argPoolId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argStorageId := args[2]

			argVote, err := cast.ToInt32E(args[3])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteBundleProposal(
				clientCtx.GetFromAddress().String(),
				argStaker,
				argPoolId,
				argStorageId,
				types.VoteType(argVote),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
