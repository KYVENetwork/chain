package cli

import (
	"strconv"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdDelegator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator [staker] [delegator]",
		Short: "Query delegator of staker",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqStaker := args[0]
			reqDelegator := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryDelegationClient(clientCtx)

			params := &types.QueryDelegatorRequest{
				Staker:    reqStaker,
				Delegator: reqDelegator,
			}

			res, err := queryClient.Delegator(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
