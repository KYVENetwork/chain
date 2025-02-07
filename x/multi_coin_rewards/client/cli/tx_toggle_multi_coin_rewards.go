package cli

import (
	"fmt"

	"github.com/KYVENetwork/chain/x/multi_coin_rewards/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdToggleMultiCoinRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle-multi-coin-rewards [enabled]",
		Short: "Broadcast message to toggle multi-coin rewards",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if args[0] != "true" && args[0] != "false" {
				return fmt.Errorf("value must be 'true' or 'false'")
			}

			msg := types.MsgToggleMultiCoinRewards{
				Creator: clientCtx.GetFromAddress().String(),
				Enabled: args[0] == "true",
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
