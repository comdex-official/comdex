package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

type AddNewAssetsRequest struct {
	BaseReq     rest.BaseReq  `json:"base_req" yaml:"base_req"`
	Title       string        `json:"title" yaml:"title"`
	Description string        `json:"description" yaml:"description"`
	Deposit     sdk.Coins     `json:"deposit" yaml:"deposit"`
	Asset       []types.Asset `json:"assets" yaml:"assets"`
}

type UpdateNewAssetRequest struct{}
type AddNewPairsRequest struct{}
type UpdateNewPairRequest struct{}

func AddNewAssetsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-assets",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}

func UpdateNewAssetProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-new-asset",
		Handler:  UpdateNewAssetRESTHandler(clientCtx),
	}
}

func UpdateNewGovTimeInAppMappingProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-gov-time-app-mapping",
		Handler:  UpdateGovTimeInAppMappingRESTHandler(clientCtx),
	}
}

func AddNewPairsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-pairs",
		Handler:  AddNewPairsRESTHandler(clientCtx),
	}
}

func AddNewAssetsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddNewAssetsRequest

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

		content := types.NewAddAssetsProposal(
			req.Title,
			req.Description,
			req.Asset)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func UpdateNewAssetRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateNewAssetRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func UpdateGovTimeInAppMappingRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateNewAssetRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func AddNewPairsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddNewPairsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func UpdateNewPairsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateNewPairRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func AddNewAppMappingProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-app-mapping",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}
func AddNewAssetMappingProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-update-new-asset-mapping",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}

func AddExtendedPairsVaultProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-pairs-vault",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}
