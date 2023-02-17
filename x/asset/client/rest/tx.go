package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

type (
	UpdateNewAssetRequest struct{}
	AddNewPairsRequest    struct{}
	UpdateNewPairRequest  struct{}
	AddNewAssetsPairs     struct{}
)

func AddNewAssetsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-assets",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}

func AddNewAssetsPairsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-assets-pairs",
		Handler:  AddNewAssetsPairsRESTHandler(clientCtx),
	}
}

func UpdateNewAssetProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-new-asset",
		Handler:  UpdateNewAssetRESTHandler(clientCtx),
	}
}

func UpdateNewPairProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-new-pair",
		Handler:  UpdateNewPairsRESTHandler(clientCtx),
	}
}

func UpdateNewGovTimeInAppProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-gov-time-app",
		Handler:  UpdateGovTimeInAppRESTHandler(clientCtx),
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
		var req UpdateNewAssetRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func AddNewAssetsPairsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddNewAssetsPairs

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
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

func UpdateGovTimeInAppRESTHandler(clientCtx client.Context) http.HandlerFunc {
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

func AddNewAppProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-app",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}

func AddNewAssetInAppProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-update-new-asset-in-app",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}

func AddExtendedPairsVaultProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-pairs-vault",
		Handler:  AddNewAssetsRESTHandler(clientCtx),
	}
}
