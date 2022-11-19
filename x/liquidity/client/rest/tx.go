package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/petrichormoney/petri/x/liquidity/types"
)

type UpdateGenericParamsRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
	AppID       uint64       `json:"app_id" yaml:"app_id"`
	Keys        []string     `json:"keys" yaml:"keys"`
	Values      []string     `json:"values" yaml:"values"`
}

type CreateNewLiquidityPairRequest struct {
	BaseReq        rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title          string       `json:"title" yaml:"title"`
	Description    string       `json:"description" yaml:"description"`
	Deposit        sdk.Coins    `json:"deposit" yaml:"deposit"`
	AppID          uint64       `json:"app_id" yaml:"app_id"`
	BaseCoinDenom  string       `json:"base_coin_denom" yaml:"base_coin_denom"`
	QuoteCoinDenom string       `json:"quote_coin_denom" yaml:"quote_coin_denom"`
}

func UpdateGenericParamsProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "liquidity-param-change",
		Handler:  UpdateGenericParamsRESTHandler(clientCtx),
	}
}

func CreateNewLiquidityPairProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "create-liquidity-pair",
		Handler:  CreateNewLiquidityPairRESTHandler(clientCtx),
	}
}

func UpdateGenericParamsRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateGenericParamsRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewUpdateGenericParamsProposal(
			req.Title,
			req.Description,
			req.AppID,
			req.Keys,
			req.Values,
		)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func CreateNewLiquidityPairRESTHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateNewLiquidityPairRequest

		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		content := types.NewCreateLiquidityPairProposal(
			req.Title,
			req.Description,
			fromAddr,
			req.AppID,
			req.BaseCoinDenom,
			req.QuoteCoinDenom,
		)
		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}
		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
