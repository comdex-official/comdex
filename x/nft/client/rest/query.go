package rest

import (
	"context"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/comdex-official/comdex/x/nft/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/%s/denoms/{%s}/supply", types.ModuleName, RestParamDenom),
		querySupply(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/collections/{%s}", types.ModuleName, RestParamDenom),
		queryCollection(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/denoms", types.ModuleName),
		queryDenoms(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/denoms/{%s}", types.ModuleName, RestParamDenom),
		queryDenom(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/owners/{%s}", types.ModuleName, RestParamOwner),
		queryOwnerNFTs(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/%s/asset/{%s}/{%s}", types.ModuleName, RestParamDenom, RestParamNFTID),
		queryNFT(cliCtx, queryRoute),
	).Methods("GET")
}

func querySupply(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		denom := strings.TrimSpace(mux.Vars(r)[RestParamDenom])
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		var owner sdk.AccAddress
		ownerStr := r.FormValue(RestParamOwner)
		if len(ownerStr) > 0 {
			ownerAddress, err := sdk.AccAddressFromBech32(strings.TrimSpace(ownerStr))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			owner = ownerAddress
		}
		params := types.NewQuerySupplyParams(denom, owner)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QuerySupply), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		out := binary.LittleEndian.Uint64(res)
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, out)
	}
}

func queryOwnerNFTs(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerStr := mux.Vars(r)[RestParamOwner]
		if len(ownerStr) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "owner should not be empty")
		}

		owner, err := sdk.AccAddressFromBech32(ownerStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		denomId := r.FormValue(RestParamDenom)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		var (
			qc = types.NewQueryClient(cliCtx)
		)

		_, page, limit, err := rest.ParseHTTPArgs(r)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		pageReq := sdkquery.PageRequest{
			Offset:     uint64((page - 1) * limit),
			Limit:      uint64(limit),
			CountTotal: true,
		}

		nfts, err := qc.OwnerNFTs(
			context.Background(),
			&types.QueryOwnerNFTsRequest{
				Owner:      owner.String(),
				DenomId:    denomId,
				Pagination: &pageReq,
			},
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, cliCtx, nfts)
	}
}

func queryCollection(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		denomId := vars[RestParamDenom]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		var (
			qc = types.NewQueryClient(cliCtx)
		)

		_, page, limit, err := rest.ParseHTTPArgs(r)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		pageReq := sdkquery.PageRequest{
			Offset:     uint64((page - 1) * limit),
			Limit:      uint64(limit),
			CountTotal: true,
		}

		collection, err := qc.Collection(
			context.Background(),
			&types.QueryCollectionRequest{
				DenomId:    denomId,
				Pagination: &pageReq,
			},
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, cliCtx, collection)
	}
}

func queryDenom(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		denom := mux.Vars(r)[RestParamDenom]
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryDenomParams(denom)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDenom), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDenoms(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		var (
			qc = types.NewQueryClient(cliCtx)
		)

		_, page, limit, err := rest.ParseHTTPArgs(r)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		pageReq := sdkquery.PageRequest{
			Offset:     uint64((page - 1) * limit),
			Limit:      uint64(limit),
			CountTotal: true,
		}
		var owner sdk.AccAddress
		if query.Get("owner") != "" {
			owner, err = sdk.AccAddressFromBech32(query.Get("owner"))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		denoms, err := qc.Denoms(
			context.Background(),
			&types.QueryDenomsRequest{
				Pagination: &pageReq,
				Owner:      owner.String(),
			},
		)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, cliCtx, denoms)
	}
}

func queryNFT(cliCtx client.Context, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		denom := vars[RestParamDenom]
		if err := types.ValidateDenomID(denom); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		nftID := vars[RestParamNFTID]
		if err := types.ValidateNFTID(nftID); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		params := types.NewQueryNFTParams(denom, nftID)
		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(
			fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryNFT), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
