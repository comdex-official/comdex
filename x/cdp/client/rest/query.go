package rest

import (
	"context"
	"github.com/gorilla/mux"
	"strconv"

	"github.com/comdex-official/comdex/x/cdp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"net/http"
)

func queryCdp(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		qc := types.NewQueryServiceClient(ctx)
		id, err := strconv.ParseUint(vars["id"], 10, 64)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := qc.QueryCDP(context.Background(),
			&types.QueryCDPRequest{
				Id: id,
			})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, res)
	}
}

func queryDeposits(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		qc := types.NewQueryServiceClient(ctx)

		res, err := qc.QueryCDPDeposits(context.Background(),
			&types.QueryCDPDepositsRequest{
				Owner:          vars["owner"],
				CollateralType: vars["collateral_type"],
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
