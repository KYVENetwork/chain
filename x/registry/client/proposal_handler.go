package client

import (
	"github.com/KYVENetwork/chain/x/registry/client/cli"
	"github.com/KYVENetwork/chain/x/registry/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var CreatePoolHandler = govclient.NewProposalHandler(cli.CmdSubmitCreatePoolProposal, rest.ProposalCreatePoolRESTHandler)
var UpdatePoolHandler = govclient.NewProposalHandler(cli.CmdSubmitUpdatePoolProposal, rest.ProposalUpdatePoolRESTHandler)
var PausePoolHandler = govclient.NewProposalHandler(cli.CmdSubmitPausePoolProposal, rest.ProposalPausePoolRESTHandler)
var UnpausePoolHandler = govclient.NewProposalHandler(cli.CmdSubmitUnpausePoolProposal, rest.ProposalUnpausePoolRESTHandler)
var SchedulePoolUpgradeHandler = govclient.NewProposalHandler(cli.CmdSubmitSchedulePoolUpgradeProposal, rest.ProposalSchedulePoolUpgradeRESTHandler)
var CancelPoolUpgradeHandler = govclient.NewProposalHandler(cli.CmdSubmitCancelPoolUpgradeProposal, rest.ProposalCancelPoolUpgradeRESTHandler)
