package cli

import (
	"context"
	"strconv"

	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/spf13/cobra"
)

func byFunder(
	clientCtx client.Context,
	queryClient types.QueryFundersClient,
	pageReq *query.PageRequest,
	address string,
) error {
	params := &types.QueryFundingsByFunderRequest{
		Pagination: pageReq,
		Address:    address,
	}

	res, err := queryClient.FundingsByFunder(context.Background(), params)
	if err != nil {
		return err
	}
	return clientCtx.PrintProto(res)
}

func byPool(
	clientCtx client.Context,
	queryClient types.QueryFundersClient,
	pageReq *query.PageRequest,
	poolId uint64,
) error {
	params := &types.QueryFundingsByPoolRequest{
		Pagination: pageReq,
		PoolId:     poolId,
	}

	res, err := queryClient.FundingsByPool(context.Background(), params)
	if err != nil {
		return err
	}
	return clientCtx.PrintProto(res)
}

func CmdListFundings() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fundings [address | pool-id]",
		Short: "list all fundings of a user or a pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reqAddressOrPool := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryFundersClient(clientCtx)

			poolId, err := strconv.ParseUint(reqAddressOrPool, 10, 64)
			if err != nil {
				return byFunder(clientCtx, queryClient, pageReq, reqAddressOrPool)
			}
			return byPool(clientCtx, queryClient, pageReq, poolId)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
