package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	"net/http"
)

type (
	WhitelistingLiquidationRequest struct{}
)

func AddNewWhitelistingLiquidationRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-liquidation-whitelisting",
		Handler:  WhitelistingLiquidationRequestRESTHandler(clientCtx),
	}
}

func WhitelistingLiquidationRequestRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WhitelistingLiquidationRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}
