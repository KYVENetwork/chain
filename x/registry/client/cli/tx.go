package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/KYVENetwork/chain/x/registry/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdFundPool())
	cmd.AddCommand(CmdDefundPool())
	cmd.AddCommand(CmdStakePool())
	cmd.AddCommand(CmdUnstakePool())
	cmd.AddCommand(CmdSubmitBundleProposal())
	cmd.AddCommand(CmdVoteProposal())
	cmd.AddCommand(CmdClaimUploaderRole())
	cmd.AddCommand(CmdDelegatePool())
	cmd.AddCommand(CmdWithdrawPool())
	cmd.AddCommand(CmdUndelegatePool())
	cmd.AddCommand(CmdRedelegatePool())
	cmd.AddCommand(CmdUpdateMetadata())
	cmd.AddCommand(CmdUpdateCommission())

	cmd.AddCommand(CmdSubmitCreatePoolProposal())
	cmd.AddCommand(CmdSubmitUpdatePoolProposal())
	cmd.AddCommand(CmdSubmitPausePoolProposal())
	cmd.AddCommand(CmdSubmitUnpausePoolProposal())
	cmd.AddCommand(CmdSubmitSchedulePoolUpgradeProposal())
	cmd.AddCommand(CmdSubmitCancelPoolUpgradeProposal())
	cmd.AddCommand(CmdSubmitResetPoolProposal())

	return cmd
}

func CmdSubmitCreatePoolProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [flags]",
		Args:  cobra.ExactArgs(11),
		Short: "Submit a proposal to create a pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			uploadInterval, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			operatingCost, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}

			maxBundleSize, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return err
			}

			minStake, err := strconv.ParseUint(args[11], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewCreatePoolProposal(title, description, args[0], args[1], args[2], args[3], uploadInterval, operatingCost, maxBundleSize, args[8], args[9], args[10], minStake)

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func CmdSubmitUpdatePoolProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool [flags]",
		Args:  cobra.ExactArgs(9),
		Short: "Submit a proposal to update a pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			uploadInterval, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return err
			}

			operatingCost, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return err
			}

			maxBundleSize, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return err
			}

			minStake, err := strconv.ParseUint(args[8], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewUpdatePoolProposal(title, description, id, args[1], args[2], args[3], args[4], uploadInterval, operatingCost, maxBundleSize, minStake)

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func CmdSubmitPausePoolProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause-pool [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a proposal to pause a pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewPausePoolProposal(title, description, id)

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func CmdSubmitUnpausePoolProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpause-pool [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a proposal to unpause a pool.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewUnpausePoolProposal(title, description, id)

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func CmdSubmitSchedulePoolUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-pool-upgrade [flags]",
		Args:  cobra.ExactArgs(5),
		Short: "Submit a proposal to schedule a pool upgrade.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			scheduledAt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			duration, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewSchedulePoolUpgradeProposal(title, description, args[0], args[1], scheduledAt, duration, args[4])

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func CmdSubmitCancelPoolUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-pool-upgrade [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a proposal to cancel a scheduled pool upgrade.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewCancelPoolUpgradeProposal(title, description, args[0])

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}

func CmdSubmitResetPoolProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-pool [flags]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a proposal to reset a pool to a specific bundle.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argBundleId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			content := types.NewResetPoolProposal(title, description, argId, argBundleId)

			isExpedited, err := cmd.Flags().GetBool(cli.FlagIsExpedited)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from, isExpedited)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "The proposal title")
	cmd.Flags().String(cli.FlagDescription, "", "The proposal description")
	cmd.Flags().Bool(cli.FlagIsExpedited, false, "If true, makes the proposal an expedited one")
	cmd.Flags().String(cli.FlagDeposit, "", "The proposal deposit")
	_ = cmd.MarkFlagRequired(cli.FlagTitle)
	_ = cmd.MarkFlagRequired(cli.FlagDescription)

	return cmd
}
