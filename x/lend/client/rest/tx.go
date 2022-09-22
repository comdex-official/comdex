package rest

import (
	"net/http"

	"github.com/comdex-official/comdex/x/asset/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

type AddNewAssetsRequest struct {
	BaseReq     rest.BaseReq  `json:"base_req" yaml:"base_req"`
	Title       string        `json:"title" yaml:"title"`
	Description string        `json:"description" yaml:"description"`
	Deposit     sdk.Coins     `json:"deposit" yaml:"deposit"`
	Asset       []types.Asset `json:"assets" yaml:"assets"`
}

type AddNewPairsRequest struct{}
type UpdateNewPairRequest struct{}
type AddPoolRequest struct{}
type AddAssetToPairRequest struct{}
type UpdateAssetToPairRequest struct{}
type AddAssetRatesStatsRequest struct{}
type AddAuctionParamsRequest struct{}

func AddNewPairsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-new-pairs",
		Handler:  AddNewPairsRESTHandler(clientCtx),
	}
}

func UpdatePairProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update-new-pair",
		Handler:  UpdateNewPairsRESTHandler(clientCtx),
	}
}

func AddPoolProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-lend-pools",
		Handler:  AddPoolRESTHandler(clientCtx),
	}
}

func AddAssetToPairProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-asset-to-pair",
		Handler:  AddAssetToPairRESTHandler(clientCtx),
	}
}

func AddNewAssetRatesStatsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-asset-rates-stats",
		Handler:  AddAssetRatesStatsRESTHandler(clientCtx),
	}
}

func AddNewAuctionParamsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-auction-params",
		Handler:  AddAuctionParamsRESTHandler(clientCtx),
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

func AddPoolRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddPoolRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func AddAssetToPairRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddAssetToPairRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func AddAssetRatesStatsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddAssetRatesStatsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func AddAuctionParamsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddAuctionParamsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}
