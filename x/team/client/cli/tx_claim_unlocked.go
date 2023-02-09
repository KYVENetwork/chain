package cli

import (
	"github.com/KYVENetwork/chain/x/team/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdClaimUnlocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-unlocked [id] [amount] [recipient]",
		Short: "Broadcast message claim-unlocked",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argAmount, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			argRecipient := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgClaimUnlocked{
				Authority: clientCtx.GetFromAddress().String(),
				Id:        argId,
				Amount:    argAmount,
				Recipient: argRecipient,
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
