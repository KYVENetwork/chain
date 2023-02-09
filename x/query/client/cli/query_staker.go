package cli

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListStakers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stakers",
		Short: "list all stakers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryStakersClient(clientCtx)

			params := &types.QueryStakersRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Stakers(context.Background(), params)
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

func CmdShowStaker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staker [address]",
		Short: "shows all necessary information for staker",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryStakersClient(clientCtx)

			params := &types.QueryStakerRequest{
				Address: args[0],
			}

			res, err := queryClient.Staker(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
