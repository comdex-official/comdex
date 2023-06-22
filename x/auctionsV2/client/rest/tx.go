package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	"net/http"
)

type (
	AddAuctionParamsRequest struct{}
)

func AddAuctionParamsRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-auction-params",
		Handler:  AuctionParamsRequestRESTHandler(clientCtx),
	}
}

func AuctionParamsRequestRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddAuctionParamsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}
