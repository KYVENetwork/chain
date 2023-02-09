package cli

import (
	"github.com/KYVENetwork/chain/x/query/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdAccountAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-assets [address]",
		Short: "Query account assets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryAccountClient(clientCtx)

			params := &types.QueryAccountAssetsRequest{
				Address: reqAddress,
			}

			res, err := queryClient.AccountAssets(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
