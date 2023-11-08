package cli

import (
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdUpdateFunder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-funder",
		Short: "Broadcast message create-funder",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			identity, _ := cmd.Flags().GetString(FlagIdentity)
			website, _ := cmd.Flags().GetString(FlagWebsite)
			contact, _ := cmd.Flags().GetString(FlagContact)
			description, _ := cmd.Flags().GetString(FlagDescription)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUpdateFunder{
				Creator:     clientCtx.GetFromAddress().String(),
				Moniker:     moniker,
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
