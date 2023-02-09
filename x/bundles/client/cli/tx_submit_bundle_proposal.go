package cli

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdSubmitBundleProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-bundle-proposal [staker] [pool_id] [storage_id] [data_size] [data_hash] [from_index] [bundle_size] [from_key] [to_key] [bundle_summary]",
		Short: "Broadcast message submit-bundle-proposal",
		Args:  cobra.ExactArgs(10),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStaker := args[0]

			argPoolId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argStorageId := args[2]

			argDataSize, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argDataHash := args[4]

			argFromIndex, err := cast.ToUint64E(args[5])
			if err != nil {
				return err
			}

			argBundleSize, err := cast.ToUint64E(args[6])
			if err != nil {
				return err
			}

			argFromKey := args[7]

			argToKey := args[8]

			argBundleSummary := args[9]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSubmitBundleProposal(
				clientCtx.GetFromAddress().String(),
				argStaker,
				argPoolId,
				argStorageId,
				argDataSize,
				argDataHash,
				argFromIndex,
				argBundleSize,
				argFromKey,
				argToKey,
				argBundleSummary,
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
