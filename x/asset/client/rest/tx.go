package rest

/*type UpdateLiquidationRatioRequest struct {
	BaseReq          rest.BaseReq           `json:"base_req" yaml:"base_req"`
	Title            string                 `json:"title" yaml:"title"`
	Description      string                 `json:"description" yaml:"description"`
	Deposit          sdk.Coins              `json:"deposit" yaml:"deposit"`
	Assets types.Asset `json:"liquidation_ration" yaml:"liquidation_ratio"`
}

func ProposalUpdateLiquidationRatioRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add-assets",
		Handler:  newUpdatePoolIncentivesHandler(clientCtx),
	}
}

func newUpdatePoolIncentivesHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateLiquidationRatioRequest

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

		content := types.NewUpdateLiquidationRatioProposal(req.Title, req.Description, req.Assets)
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
*/
