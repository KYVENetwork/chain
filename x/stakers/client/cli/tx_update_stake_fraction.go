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

func CmdUpdateStakeFraction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-stake-fraction [pool_id] [stake_fraction]",
		Short: "Broadcast message update-stake-fraction",
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

			argStakeFraction, err := math.LegacyNewDecFromStr(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgUpdateStakeFraction{
				Creator:       clientCtx.GetFromAddress().String(),
				PoolId:        argPoolId,
				StakeFraction: argStakeFraction,
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
