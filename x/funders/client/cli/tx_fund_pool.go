package cli

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdFundPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fund-pool [id] [amount] [amount_per_bundle]",
		Short: "Broadcast message fund-pool",
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
			argAmountPerBundle, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgFundPool{
				Creator:         clientCtx.GetFromAddress().String(),
				PoolId:          argId,
				Amount:          argAmount,
				AmountPerBundle: argAmountPerBundle,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
