package cli

import (
	"context"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListFunders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "funders",
		Short: "list all funders",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reqSearch := ""
			if len(args) >= 1 {
				reqSearch = args[0]
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryFundersClient(clientCtx)

			params := &types.QueryFundersRequest{
				Pagination: pageReq,
				Search:     reqSearch,
			}

			res, err := queryClient.Funders(context.Background(), params)
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

func CmdShowFunder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "funder [address]",
		Short: "shows a funder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reqAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryFundersClient(clientCtx)

			params := &types.QueryFunderRequest{
				Address: reqAddress,
			}

			res, err := queryClient.Funder(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
