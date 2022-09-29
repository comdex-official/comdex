package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
)

type (
	NewCmdLookupTableParams       struct{}
	NewCmdAuctionToAppLookupTable struct{}
)

func NewCmdLookupTableParamsRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "collector-lookup-params",
		Handler:  LookupTableParamsRESTHandler(clientCtx),
	}
}

func LookupTableParamsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req NewCmdLookupTableParams

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}

func NewCmdAuctionTableAppRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "app-auction-lookup",
		Handler:  AuctionLookupTableAppRESTHandler(clientCtx),
	}
}

func AuctionLookupTableAppRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req NewCmdAuctionToAppLookupTable

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}
