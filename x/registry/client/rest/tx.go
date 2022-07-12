package rest

import (
	"net/http"

	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type CreatePoolRequest struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title          string       `json:"title" yaml:"title"`
	Description    string       `json:"description" yaml:"description"`
	IsExpedited    bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
	Name           string       `json:"name" yaml:"name"`
	Runtime        string       `json:"runtime" yaml:"runtime"`
	Logo           string       `json:"logo" yaml:"logo"`
	Config         string       `json:"config" yaml:"config"`
	UploadInterval uint64       `json:"uploadInterval" yaml:"uploadInterval"`
	OperatingCost  uint64       `json:"operatingCost" yaml:"operatingCost"`
	MaxBundleSize  uint64       `json:"maxBundleSize" yaml:"maxBundleSize"`
	Version        string       `json:"version" yaml:"version"`
	Binaries       string       `json:"binaries" yaml:"binaries"`
	StartKey       string       `json:"startKey" yaml:"startKey"`
	MinStake       uint64       `json:"minStake" yaml:"minStake"`
}

func ProposalCreatePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "create-pool",
		Handler:  newCreatePoolHandler(clientCtx),
	}
}

func newCreatePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewCreatePoolProposal(req.Title, req.Description, req.Name, req.Runtime, req.Logo, req.Config, req.UploadInterval, req.OperatingCost, req.MaxBundleSize, req.Version, req.Binaries, req.StartKey, req.MinStake)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type UpdatePoolRequest struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title          string       `json:"title" yaml:"title"`
	Description    string       `json:"description" yaml:"description"`
	IsExpedited    bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id             uint64       `json:"id" yaml:"id"`
	Name           string       `json:"name" yaml:"name"`
	Runtime        string       `json:"runtime" yaml:"runtime"`
	Logo           string       `json:"logo" yaml:"logo"`
	Config         string       `json:"config" yaml:"config"`
	UploadInterval uint64       `json:"uploadInterval" yaml:"uploadInterval"`
	OperatingCost  uint64       `json:"operatingCost" yaml:"operatingCost"`
	MaxBundleSize  uint64       `json:"maxBundleSize" yaml:"maxBundleSize"`
	MinStake       uint64       `json:"minStake" yaml:"minStake"`
}

func ProposalUpdatePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-pool",
		Handler:  newUpdatePoolHandler(clientCtx),
	}
}

func newUpdatePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdatePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewUpdatePoolProposal(req.Title, req.Description, req.Id, req.Name, req.Runtime, req.Logo, req.Config, req.UploadInterval, req.OperatingCost, req.MaxBundleSize, req.MinStake)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type PausePoolRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	IsExpedited bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id          uint64       `json:"id" yaml:"id"`
}

func ProposalPausePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "pause-pool",
		Handler:  newPausePoolHandler(clientCtx),
	}
}

func newPausePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PausePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewPausePoolProposal(req.Title, req.Description, req.Id)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type UnpausePoolRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	IsExpedited bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id          uint64       `json:"id" yaml:"id"`
}

func ProposalUnpausePoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "unpause-pool",
		Handler:  newUnpausePoolHandler(clientCtx),
	}
}

func newUnpausePoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnpausePoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewUnpausePoolProposal(req.Title, req.Description, req.Id)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type SchedulePoolUpgradeRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	IsExpedited bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Runtime     string       `json:"runtime" yaml:"runtime"`
	Version     string       `json:"version" yaml:"version"`
	ScheduledAt uint64       `json:"scheduled_at" yaml:"scheduled_at"`
	Duration    uint64       `json:"duration" yaml:"duration"`
	Binaries    string       `json:"binaries" yaml:"binaries"`
}

func ProposalSchedulePoolUpgradeRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "schedule-pool-upgrade",
		Handler:  newSchedulePoolUpgradeHandler(clientCtx),
	}
}

func newSchedulePoolUpgradeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SchedulePoolUpgradeRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewSchedulePoolUpgradeProposal(req.Title, req.Description, req.Runtime, req.Version, req.ScheduledAt, req.Duration, req.Binaries)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type CancelPoolUpgradeRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	IsExpedited bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Runtime     string       `json:"runtime" yaml:"runtime"`
}

func ProposalCancelPoolUpgradeRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "cancel-pool-upgrade",
		Handler:  newCancelPoolUpgradeHandler(clientCtx),
	}
}

func newCancelPoolUpgradeHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CancelPoolUpgradeRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewCancelPoolUpgradeProposal(req.Title, req.Description, req.Runtime)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

type ResetPoolRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	IsExpedited bool         `json:"is_expedited" yaml:"is_expedited"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	Id     uint64       `json:"id" yaml:"id"`
	BundleId     uint64       `json:"bundleId" yaml:"bundleId"`
}

func ProposalResetPoolRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "cancel-pool-upgrade",
		Handler:  newCancelPoolUpgradeHandler(clientCtx),
	}
}

func newResetPoolHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ResetPoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewResetPoolProposal(req.Title, req.Description, req.Id, req.BundleId)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr, req.IsExpedited)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
