package cli

import (
	"github.com/KYVENetwork/chain/x/bundles/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdSkipUploaderRole() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skip-uploader-role [staker] [pool_id] [from_index]",
		Short: "Broadcast message skip-uploader-role",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argStaker := args[0]

			argPoolId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argFromIndex, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSkipUploaderRole(
				clientCtx.GetFromAddress().String(),
				argStaker,
				argPoolId,
				argFromIndex,
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
