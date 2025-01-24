package cli

import (
	"encoding/json"
	"os"

	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdSetMultiCoinRefundPolicy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-multi-coin-refund-policy [path-to-json-file]",
		Short: "Broadcast message to update the refund policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			file, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			policy := &types.MultiCoinRefundPolicy{}
			err = json.Unmarshal(file, policy)
			if err != nil {
				return err
			}

			msg := types.MsgSetMultiCoinRewardsRefundPolicy{
				Creator: clientCtx.GetFromAddress().String(),
				Policy:  policy,
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
