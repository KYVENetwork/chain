package cli

import (
	"context"

	"github.com/spf13/cast"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListFinalizedBundles() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finalized-bundles [pool_id]",
		Short: "list all finalized bundles of pool given by pool_id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			poolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryBundlesClient(clientCtx)

			params := &types.QueryFinalizedBundlesRequest{
				PoolId:     poolId,
				Pagination: pageReq,
			}

			res, err := queryClient.FinalizedBundlesQuery(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowFinalizedBundle() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finalized-bundle [pool_id] [bundle-id]",
		Short: "show the finalized bundle given by pool_id and bundle_id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			poolId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			bundleId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryBundlesClient(clientCtx)

			params := &types.QueryFinalizedBundleRequest{
				PoolId: poolId,
				Id:     bundleId,
			}

			res, err := queryClient.FinalizedBundleQuery(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
