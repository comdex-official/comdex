package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterHandlers(cliCtx client.Context, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r, queryRoute)
}

const (
	RestParamDenom = "denom"
	RestParamNFTID = "id"
	RestParamOwner = "owner"
)

type createDenomReq struct {
	BaseReq     rest.BaseReq   `json:"base_req"`
	Sender      sdk.AccAddress `json:"sender"`
	Symbol      string         `json:"symbol"`
	Name        string         `json:"name"`
	Schema      string         `json:"schema"`
	Description string         `json:"description"`
	PreviewURI  string         `json:"preview_uri"`
}

type updateDenomReq struct {
	BaseReq     rest.BaseReq   `json:"base_req"`
	Sender      sdk.AccAddress `json:"sender"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	PreviewURI  string         `json:"preview_uri"`
}

type transferDenomReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Sender    sdk.AccAddress `json:"sender"`
	Recipient string         `json:"recipient"`
}

type mintNFTReq struct {
	BaseReq      rest.BaseReq   `json:"base_req"`
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    sdk.AccAddress `json:"recipient"`
	Denom        string         `json:"denom"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	MediaURI     string         `json:"media_uri"`
	PreviewURI   string         `json:"preview_uri"`
	Data         string         `json:"data"`
	Transferable string         `json:"transferable"`
	Extensible   string         `json:"extensible"`
	Nsfw         string         `json:"nsfw"`
	RoyaltyShare sdk.Dec        `json:"royalty_share"`
}

type transferNFTReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Sender    sdk.AccAddress `json:"sender"`
	Recipient string         `json:"recipient"`
}

type burnNFTReq struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Sender  sdk.AccAddress `json:"sender"`
}
