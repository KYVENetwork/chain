package cli

import (
	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdUpdateCommission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-commission [pool_id] [commission]",
		Short: "Broadcast message update-commission",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argCommission, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgUpdateCommission{
				Creator:    clientCtx.GetFromAddress().String(),
				PoolId:     argPoolId,
				Commission: argCommission,
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
