package cli

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateFunder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-funder [moniker]",
		Short: "Broadcast message create-funder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argMoniker := args[0]
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			contact, _ := cmd.Flags().GetString(FlagContact)
			description, _ := cmd.Flags().GetString(FlagDescription)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgCreateFunder{
				Creator:     clientCtx.GetFromAddress().String(),
				Moniker:     argMoniker,
				Identity:    identity,
				Website:     website,
				Contact:     contact,
				Description: description,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(flagSetFunderCreate())
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
