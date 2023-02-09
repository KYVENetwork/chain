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

func CmdCanVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "can-vote [pool_id] [storage-id] [voter]",
		Short: "Query if the current voter can vote on the current proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqStorageId := args[1]
			reqVoter := args[2]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryBundlesClient(clientCtx)

			params := &types.QueryCanVoteRequest{
				PoolId:    reqId,
				StorageId: reqStorageId,
				Voter:     reqVoter,
			}

			res, err := queryClient.CanVote(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
