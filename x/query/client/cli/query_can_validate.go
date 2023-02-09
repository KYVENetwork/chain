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

func CmdCanValidate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "can-validate [pool_id] [valaddress]",
		Short: "Query if current valaddress can vote in pool",
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

			params := &types.QueryCanValidateRequest{
				PoolId:     reqId,
				Valaddress: args[1],
			}

			res, err := queryClient.CanValidate(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
