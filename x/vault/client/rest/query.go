package rest

import (
	"context"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/comdex-official/comdex/x/vault/types"
)

func queryVault(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		qc := types.NewQueryClient(ctx)
		idParam := vars["id"]

		id := idParam

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
