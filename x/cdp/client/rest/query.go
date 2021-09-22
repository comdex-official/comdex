package rest

import (
	"context"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/comdex-official/comdex/x/cdp/types"
)

func queryCdp(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		qc := types.NewQueryServiceClient(ctx)
		owner := vars["owner"]
		collateralType := vars["collateral_type"]

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
