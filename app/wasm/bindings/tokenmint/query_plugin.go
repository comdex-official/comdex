package tokenmint

import (
	"encoding/json"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomQuerier(tokenMintKeeper *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery TokenMintQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "tokenMint query")
		}
		if contractQuery.MintedToken != nil {
			App_Id := contractQuery.MintedToken.App_Id
			Asset_Id := contractQuery.MintedToken.Asset_Id
			MintedToken, _ := tokenMintKeeper.GetTokenMint(ctx, App_Id, Asset_Id)
			res := MintedTokenResponse{
				MintedTokens: MintedToken,
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "tokenMint query response")
			}
			return bz, nil
		}

		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown tokenMint query variant"}
	}

}
