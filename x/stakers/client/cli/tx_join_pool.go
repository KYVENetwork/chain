package cli

import (
	"github.com/KYVENetwork/chain/x/stakers/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdJoinPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join-pool [pool_id] [valaddress] [amount]",
		Short: "Broadcast message join-pool",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argValaddress := args[1]

			argAmount, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgJoinPool{
				Creator:    clientCtx.GetFromAddress().String(),
				PoolId:     argPoolId,
				Valaddress: argValaddress,
				Amount:     argAmount,
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
