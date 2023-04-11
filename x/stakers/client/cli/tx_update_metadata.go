package cli

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdUpdateMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-metadata [moniker] [website] [identity] [security_contact] [details]",
		Short: "Broadcast message update-metadata",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgUpdateMetadata{
				Creator:         clientCtx.GetFromAddress().String(),
				Moniker:         args[0],
				Website:         args[1],
				Identity:        args[2],
				SecurityContact: args[3],
				Details:         args[4],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
