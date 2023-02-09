package cli

import (
	"github.com/KYVENetwork/chain/x/team/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateTeamVestingAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [total_allocation] [commencement]",
		Short: "Broadcast message create-team-vesting-account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argAllocation, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argCommencementTimeStamp, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgCreateTeamVestingAccount{
				Authority:       clientCtx.GetFromAddress().String(),
				TotalAllocation: argAllocation,
				Commencement:    argCommencementTimeStamp,
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
