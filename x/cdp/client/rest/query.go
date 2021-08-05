package rest

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func queryCdp(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		owner := vars["owner"]
		collateralType := vars["collateral_type"]

		qc := types.NewQueryServiceClient(ctx)

		res, err := qc.QueryCDP(context.Background(),
			&types.QueryCDPRequest{
				Owner:          owner,
				CollateralType: collateralType,
			})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, res)
	}
}

func queryParams(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		qc := types.NewQueryServiceClient(ctx)

		res, err := qc.QueryParams(context.Background(), &types.QueryParamsRequest{})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, res)
	}
}
