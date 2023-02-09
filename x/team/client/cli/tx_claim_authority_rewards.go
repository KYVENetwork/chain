package cli

import (
	"github.com/KYVENetwork/chain/x/team/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdClaimAuthorityRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-authority-rewards [amount] [recipient]",
		Short: "Broadcast message claim-authority-rewards",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAmount, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argRecipient := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgClaimAuthorityRewards{
				Authority: clientCtx.GetFromAddress().String(),
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
