package cli

import (
	"strconv"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCurrentVoteStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-vote-status [pool_id]",
		Short: "Query current vote tally of pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryBundlesClient(clientCtx)

			params := &types.QueryCurrentVoteStatusRequest{
				PoolId: reqId,
			}

			res, err := queryClient.CurrentVoteStatus(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
