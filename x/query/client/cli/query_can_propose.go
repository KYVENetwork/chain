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

func CmdCanPropose() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "can-propose [pool-id] [proposer] [from-height]",
		Short: "Query if node can propose next bundle",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqProposer := args[1]
			reqFromIndex, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryBundlesClient(clientCtx)

			params := &types.QueryCanProposeRequest{
				PoolId:    reqId,
				Proposer:  reqProposer,
				FromIndex: reqFromIndex,
			}

			res, err := queryClient.CanPropose(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
