package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/petrichormoney/petri/x/vault/types"
)

func queryVault(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		qc := types.NewQueryClient(ctx)
		idParam := vars["id"]

		id, _ := strconv.ParseUint(idParam, 10, 64)

		res, err := qc.QueryVault(context.Background(),
			&types.QueryVaultRequest{
				Id: id,
			})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, ctx, res)
	}
}
